//
// Copyright 2018 Digital Transaction Limited. All Rights Reserved.
//
// Authors: Yang SONG, Eric Ma
//

package parallelcore_client_sdk_go

import (
	"context"
	"crypto/tls"
	"fmt"
	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"math/rand"
	"strings"
	"time"
)

// customCredential, which stores the JWT
type customCredential struct {
	token string
}

func (t customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (t customCredential) RequireTransportSecurity() bool {
	return true
}

type Client struct {
	conn       *grpc.ClientConn
	ctx        context.Context
	cancel     context.CancelFunc
	grpcClient pb.RequestHandlerClient
	token      string
}

func openOne(endpoint string, certPath string, token string) (*Client, error) {
	var grpcOpts []grpc.DialOption

	//if certPath == "" {
	//	fmt.Println("Using system default root certificates.")
	//}

	if token == "" {

		//fmt.Println("Setup a connection to fetch token.\n")
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))

	} else {

		//fmt.Println("Setup a connection with token.\n")
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(customCredential{
			token: token,
		}))

	}

	conn, err := grpc.Dial(endpoint, grpcOpts...)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Failed to dail. %v", err)
	}

	grpcClient := pb.NewRequestHandlerClient(conn)
	ctx, cancel := context.WithCancel(context.Background())

	response, err := grpcClient.Ping(ctx, &pb.Request{Payload: []byte("")})
	if err != nil {
		cancel()
		conn.Close()
		return nil, fmt.Errorf("Failed to Ping. %v", err)
	}
	if string(response.Payload) != "Pong" {
		cancel()
		conn.Close()
		return nil, fmt.Errorf("Ping wrong response value. %v", response.Payload)
	}

	return &Client{conn, ctx, cancel, grpcClient, token}, nil
}

// certPath: if empty, use the system certificate, otherwise, use the certificate provided in the file in certPath
// No TLS support: (TODO) use environment variable to control possibly..

func OpenAny(endpointSpecs string, clientId string, credential string, certPath string) (*Client, error) {
	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)

	// do until endpoints is empty:
	//  randomly select endpoint from endpoints
	//  try to openOne(endpoint) -> c
	//  if ok
	//    return c
	//  remove endpoint from endpoints
	// return error

	var lastError error
	endpoints := strings.Split(endpointSpecs, " ")
	for len(endpoints) != 0 {
		randPick := randGen.Intn(len(endpoints))
		endpoint := endpoints[randPick]

		client, err := openOne(endpoint, certPath, "")
		if err != nil {
			lastError = err
			endpoints = append(endpoints[:randPick], endpoints[randPick+1:]...)
			continue
		}

		// Successfully setup a connection, fetch the token
		tokenBytes, err := client.auth([]byte(clientId), []byte(credential))
		client.Close()
		if err != nil {
			return nil, err
		}

		//fmt.Println("token: ", string(tokenBytes), "\n")
		client, err = openOne(endpoint, certPath, string(tokenBytes))
		if err != nil {
			return nil, fmt.Errorf("CLIENT: Failed to open a connection. %v", err)
		}
		return client, nil
	}

	return nil, fmt.Errorf("CLIENT: Failed to open any client. Last error: %v", lastError)
}

// certPath: check OpenAny usage

func OpenMany(endpointSpecs string, clientId string, credential string, certPath string) ([]*Client, error) {
	endpoints := strings.Split(endpointSpecs, " ")
	lastError := fmt.Errorf("")
	clients := make([]*Client, 0)
	// these clients can use the same token
	token := ""
	for _, endpoint := range endpoints {
		client, err := openOne(endpoint, certPath, token)
		if err != nil {
			lastError = fmt.Errorf("%v%v\n", lastError, err)
			continue
		}
		if token == "" {
			// Successfully setup the first connection, fetch the token
			tokenBytes, err := client.auth([]byte(clientId), []byte(credential))
			client.Close()
			if err != nil {
				return nil, err
			}

			// Saved the token, later, all connections will setup with that token
			token = string(tokenBytes)
			//fmt.Println("token: ", token, "\n")
			client, err = openOne(endpoint, certPath, token)
			if err != nil {
				return nil, fmt.Errorf("CLIENT: Failed to open a connection. %v", err)
			}
		}
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		// All attemps to connect are failed.
		return nil, fmt.Errorf("CLIENT: Failed to open all clients. Last error: %v", lastError)
	} else {
		// All good or some attemps to connect failed. Opened all clients or opened some clients.
		return clients, lastError
	}
}

func CloseMany(clients []*Client) {
	for _, client := range clients {
		client.Close()
	}
}

func (client *Client) Close() {
	client.cancel()
	client.conn.Close()
}

func (client *Client) GetToken() string {
	return client.token
}

func (client *Client) invoke(in []byte) ([]byte, error) {
	response, err := client.grpcClient.Invoke(client.ctx, &pb.Request{Payload: in})
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	return response.Payload, nil
}

func (client *Client) Invoke(smartcontract_spec string, args []byte) ([]byte, error) {
	newBytes := append([]byte(smartcontract_spec), []byte(" ")...)
	newBytes = append(newBytes, args...)
	return client.invoke(newBytes)
}

func (client *Client) auth(clientId []byte, credential []byte) ([]byte, error) {
	response, err := client.grpcClient.Auth(client.ctx, &pb.AuthRequest{ClientId: clientId, Credential: credential})
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	return response.Payload, nil
}

func (client *Client) SysMan(in []byte) ([]byte, error) {
	response, err := client.grpcClient.SysMan(client.ctx, &pb.Request{Payload: in})
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	return response.Payload, nil
}
