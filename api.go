//
// Copyright 2019 Digital Transaction Limited. All Rights Reserved.
//
// Authors: Yang SONG, Eric Ma, Ray Chan
//

package parallelcore_client_sdk_go

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// customCredential, which stores the JWT
type customCredential struct {
	token string
}

// UserManData for parsing JSON
type UserManData struct {
	Action string `json:"action"`
	Data   []byte `json:"data"`
}

type BlockData struct {
	ChainId string `json: "chainId"`
	BlockId string `json: "blockId"`
}

type ClientData struct {
	ID         string `json:"clientId"`
	Credential string `json:"clientCredential"`
	Roles      string `json:"clientRoles"`
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
	conn            *grpc.ClientConn
	grpcClient      pb.RequestHandlerClient
	endpointSpecs   string
	certPath        string
	token           string
	expireTimestamp int64
}

func OpenAnyWithCert(endpointSpecs string, clientId string, credential string) (*Client, error) {
	tlsCertPath := os.Getenv("PCORE_CERT_PATH")
	if tlsCertPath == "" {
		return OpenAny(endpointSpecs, clientId, credential, "")
	} else {
		return OpenAny(endpointSpecs, clientId, credential, tlsCertPath)
	}
}

func OpenAnyByTokenWithCert(endpointSpecs string, token string, expireTimestamp int64) (*Client, error) {
	tlsCertPath := os.Getenv("PCORE_CERT_PATH")
	if tlsCertPath == "" {
		return OpenAnyByToken(endpointSpecs, token, expireTimestamp, "")
	} else {
		return OpenAnyByToken(endpointSpecs, token, expireTimestamp, tlsCertPath)
	}
}

func OpenManyWithCert(endpointSpecs string, clientId string, credential string) ([]*Client, error) {
	tlsCertPath := os.Getenv("PCORE_CERT_PATH")
	if tlsCertPath == "" {
		return OpenMany(endpointSpecs, clientId, credential, "")
	} else {
		return OpenMany(endpointSpecs, clientId, credential, tlsCertPath)
	}
}

func OpenManyByTokenWithCert(endpointSpecs string, token string, expireTimestamp int64) ([]*Client, error) {
	tlsCertPath := os.Getenv("PCORE_CERT_PATH")
	if tlsCertPath == "" {
		return OpenManyByToken(endpointSpecs, token, expireTimestamp, "")
	} else {
		return OpenManyByToken(endpointSpecs, token, expireTimestamp, tlsCertPath)
	}
}

func openOne(endpoint string, certPath string, token string) (*Client, error) {
	var grpcOpts []grpc.DialOption

	if token == "" {

		//fmt.Println("Setup a connection to fetch token.\n")
		if certPath == "" {
			grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		} else {
			creds, err := credentials.NewClientTLSFromFile(certPath, "")
			if err != nil {
				return nil, err
			}
			grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(creds))
		}

	} else {

		//fmt.Println("Setup a connection with token.\n")
		if certPath == "" {
			grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
		} else {
			creds, err := credentials.NewClientTLSFromFile(certPath, "")
			if err != nil {
				return nil, err
			}
			grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(creds))
		}
		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(customCredential{
			token: token,
		}))

	}

	conn, err := grpc.Dial(endpoint, grpcOpts...)
	if err != nil {
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

	cancel()
	return &Client{conn, grpcClient, "", certPath, token, 0}, nil
}

// certPath: if empty, use the system certificate, otherwise, use the certificate provided in the file in certPath

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

		// Successfully setup a connection, fetch the token and expireTimestamp
		returnBytes, err := client.auth([]byte(clientId), []byte(credential))
		client.Close()
		if err != nil {
			return nil, err
		}

		// Parse returnBytes
		tmpArray := strings.Split(string(returnBytes), " ")
		if len(tmpArray) != 2 {
			return nil, fmt.Errorf("SYSTEM: Wrong format message from auth().")
		}
		token := tmpArray[0]
		expireTimestamp, err := strconv.ParseInt(tmpArray[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("SYSTEM: Wrong format of expireTimestamp. %v", err)
		}

		client, err = openOne(endpoint, certPath, token)
		if err != nil {
			return nil, fmt.Errorf("CLIENT: Failed to open a connection. %v", err)
		}
		client.expireTimestamp = expireTimestamp
		client.endpointSpecs = endpointSpecs

		return client, nil
	}

	return nil, fmt.Errorf("CLIENT: Failed to open any client. Last error: %v", lastError)
}

// certPath: check OpenAny usage
// expireTimestamp: could be set to -1, if user do not want to track when to renew the client

func OpenAnyByToken(endpointSpecs string, token string, expireTimestamp int64, certPath string) (*Client, error) {
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

		client, err := openOne(endpoint, certPath, token)
		if err != nil {
			return nil, fmt.Errorf("CLIENT: Failed to open a connection using token. %v", err)
		}
		client.expireTimestamp = expireTimestamp
		client.endpointSpecs = endpointSpecs

		return client, nil
	}

	return nil, fmt.Errorf("CLIENT: Failed to open any client using token. Last error: %v", lastError)
}

// certPath: check OpenAny usage

func OpenMany(endpointSpecs string, clientId string, credential string, certPath string) ([]*Client, error) {
	endpoints := strings.Split(endpointSpecs, " ")
	lastError := fmt.Errorf("")
	clients := make([]*Client, 0)
	// these clients can use the same token
	token := ""
	var expireTimestamp int64
	for _, endpoint := range endpoints {
		// Create a connection first.
		// If the connection is not connect with token, then will run next block to fetch JWT and expireTimestamp
		client, err := openOne(endpoint, certPath, token)
		if err != nil {
			lastError = fmt.Errorf("%v%v\n", lastError, err)
			continue
		}
		if token == "" {
			// Successfully setup the first connection, fetch the token
			returnBytes, err := client.auth([]byte(clientId), []byte(credential))
			client.Close()
			if err != nil {
				return nil, err
			}

			// Parse returnBytes
			tmpArray := strings.Split(string(returnBytes), " ")
			if len(tmpArray) != 2 {
				return nil, fmt.Errorf("SYSTEM: Wrong format message from auth().")
			}
			// Saved the token and expireTimestamp. Later, all connections will setup with that token
			token = tmpArray[0]
			expireTimestamp, err = strconv.ParseInt(tmpArray[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("SYSTEM: Wrong format of expireTimestamp. %v", err)
			}

			client, err = openOne(endpoint, certPath, token)
			if err != nil {
				return nil, fmt.Errorf("CLIENT: Failed to open a connection. %v", err)
			}
		}
		client.expireTimestamp = expireTimestamp
		client.endpointSpecs = endpointSpecs
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		// All attemps to connect are failed.
		return nil, fmt.Errorf("CLIENT: Failed to open all clients. Last error: %v", lastError)
	} else {
		if lastError.Error() == "" {
			return clients, nil
		}
		// All good or some attemps to connect failed. Opened all clients or opened some clients.
		return clients, lastError
	}
}

func OpenManyByToken(endpointSpecs string, token string, expireTimestamp int64, certPath string) ([]*Client, error) {
	endpoints := strings.Split(endpointSpecs, " ")
	lastError := fmt.Errorf("")
	clients := make([]*Client, 0)
	for _, endpoint := range endpoints {
		client, err := openOne(endpoint, certPath, token)
		if err != nil {
			lastError = fmt.Errorf("%v%v\n", lastError, err)
			continue
		}
		client.expireTimestamp = expireTimestamp
		client.endpointSpecs = endpointSpecs
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		// All attemps to connect are failed.
		return nil, fmt.Errorf("CLIENT: Failed to open all clients using token. Last error: %v", lastError)
	} else {
		if lastError.Error() == "" {
			return clients, nil
		}
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
	client.conn.Close()
}

func (client *Client) GetToken() string {
	return client.token
}

func (client *Client) GetTokenExpTime() int64 {
	return client.expireTimestamp
}

func (client *Client) invoke(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.Invoke(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) Invoke(smartcontract_spec string, args []byte) ([]byte, error) {
	newBytes := append([]byte(smartcontract_spec), []byte(" ")...)
	newBytes = append(newBytes, args...)
	return client.invoke(newBytes)
}

func (client *Client) identifiedInvoke(in []byte) ([]byte, string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.IdentifiedInvoke(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, "", fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, "", fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, string(response.CommittedId), nil
}

func (client *Client) IdentifiedInvoke(smartcontract_spec string, args []byte) ([]byte, string, error) {
	newBytes := append([]byte(smartcontract_spec), []byte(" ")...)
	newBytes = append(newBytes, args...)
	return client.identifiedInvoke(newBytes)
}

func (client *Client) auth(clientId []byte, credential []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.Auth(ctx, &pb.AuthRequest{ClientId: clientId, Credential: credential})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) SysMan(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.SysMan(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) RegisterSmartContract(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.RegisterSmartContract(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) CreateDomain(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.CreateDomain(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) ListDomain(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.ListDomain(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) ListManagedDomains(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.ListManagedDomains(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) GrantDomainAdmin(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.GrantDomainAdmin(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) RevokeDomainAdmin(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.RevokeDomainAdmin(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) CreateClient(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.CreateClient(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) UpdateClient(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.UpdateClient(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) ListClient(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.ListClient(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) ListClients(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.ListClients(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) RemoveClient(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.RemoveClient(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) GrantAccess(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.GrantAccess(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) RevokeAccess(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.RevokeAccess(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) ListSmartContract(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.ListSmartContract(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) ListSmartContracts(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.ListSmartContracts(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) UpdateSelfCredential(clientId, credential string) ([]byte, error) {

	userManData := UserManData{}
	userManData.Action = "updateSelfCredential"
	clientData := ClientData{}
	clientData.ID = clientId
	clientData.Credential = credential
	// roles is empty!
	jsonDataClient, err := json.Marshal(clientData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	userManData.Data = jsonDataClient

	jsonData, err := json.Marshal(userManData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	return client.userMan(jsonData)
}

func (client *Client) ListInvokableSC() ([]byte, error) {
	userManData := UserManData{}
	userManData.Action = "ListInvokableSC"
	opts, err := json.Marshal(userManData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}
	return client.userMan(opts)
}

func (client *Client) GetBlockchainSummaryJson() ([]byte, error) {
	userManData := UserManData{}
	userManData.Action = "GetBlockchainSummaryJson"
	jsonData, err := json.Marshal(userManData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}
	return client.userMan(jsonData)
}

func (client *Client) GetBlockDetailsJson(chainId, blockId string) ([]byte, error) {
	userManData := UserManData{}
	userManData.Action = "GetBlockDetailsJson"
	blockData := BlockData{}
	blockData.ChainId = chainId
	blockData.BlockId = blockId

	jsonDataBlock, err := json.Marshal(blockData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	userManData.Data = jsonDataBlock

	jsonData, err := json.Marshal(userManData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}
	return client.userMan(jsonData)
	// return client.userMan([]byte("GetBlockDetailsJson " + chainId + " " + blockId))
}

func (client *Client) CalculateBlockHash(chainId, blockId string) ([]byte, error) {
	userManData := UserManData{}
	userManData.Action = "CalculateBlockHash"
	blockData := BlockData{}
	blockData.ChainId = chainId
	blockData.BlockId = blockId

	jsonDataBlock, err := json.Marshal(blockData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	userManData.Data = jsonDataBlock

	jsonData, err := json.Marshal(userManData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	return client.userMan(jsonData)
	// return client.userMan([]byte("CalculateBlockHash " + chainId + " " + blockId))
}

func (client *Client) GetSmartContractTransactionJson(transactionId string) ([]byte, error) {
	userManData := UserManData{}
	userManData.Action = "GetSmartContractTransactionJson"
	userManData.Data = []byte(transactionId)
	jsonData, err := json.Marshal(userManData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}
	return client.userMan(jsonData)
	// return client.userMan([]byte("GetSmartContractTransactionJson " + transactionId))
}

func (client *Client) GetSmartContractTransactionMetadataJson(transactionId string) ([]byte, error) {
	userManData := UserManData{}
	userManData.Action = "GetSmartContractTransactionMetadataJson"
	userManData.Data = []byte(transactionId)
	jsonData, err := json.Marshal(userManData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}
	return client.userMan(jsonData)
	// return client.userMan([]byte("GetSmartContractTransactionMetadataJson " + transactionId))
}

func (client *Client) ListLatestTransactions(count int) ([]byte, error) {
	userManData := UserManData{}
	userManData.Action = "ListLatestTransactions"

	// integer check in cli main.go
	jsonCount, err := json.Marshal(count)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	userManData.Data = jsonCount

	jsonData, err := json.Marshal(userManData)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: %s\n", err.Error())
	}
	return client.userMan(jsonData)
	// return client.userMan([]byte("ListLatestTransactions " + strconv.Itoa(count)))
}

func (client *Client) userMan(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.UserMan(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return nil, fmt.Errorf("%v", string(response.Error))
	}
	cancel()
	return response.Payload, nil
}

func (client *Client) renewToken() (string, int64, error) {
	// Fetch new JWT and expireTimestamp
	ctx, cancel := context.WithCancel(context.Background())
	response, err := client.grpcClient.Renew(ctx, &pb.Request{Payload: []byte("")})
	if err != nil {
		cancel()
		return "", 0, fmt.Errorf("CLIENT: %v", err)
	}
	if len(response.Error) != 0 {
		cancel()
		return "", 0, fmt.Errorf("%v", string(response.Error))
	}

	// Parse returnBytes
	tmpArray := strings.Split(string(response.Payload), " ")
	if len(tmpArray) != 2 {
		cancel()
		return "", 0, fmt.Errorf("SYSTEM: Wrong format message from RenewToken().")
	}
	token := tmpArray[0]
	expireTimestamp, err := strconv.ParseInt(tmpArray[1], 10, 64)
	if err != nil {
		cancel()
		return "", 0, fmt.Errorf("SYSTEM: Wrong format of expireTimestamp. %v", err)
	}

	cancel()
	return token, expireTimestamp, nil
}

func (client *Client) Renew() error {

	token, expireTimestamp, err := client.renewToken()
	if err != nil {
		return err
	}

	client.Close()

	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)

	endpoints := strings.Split(client.endpointSpecs, " ")
	var lastError error
	for len(endpoints) != 0 {
		randPick := randGen.Intn(len(endpoints))
		endpoint := endpoints[randPick]

		newClient, err := openOne(endpoint, client.certPath, token)
		if err != nil {
			lastError = err
			endpoints = append(endpoints[:randPick], endpoints[randPick+1:]...)
			continue
		}
		newClient.expireTimestamp = expireTimestamp
		newClient.endpointSpecs = client.endpointSpecs

		*client = *newClient
		return nil
	}
	return fmt.Errorf("CLIENT: Failed to renew a connection. %v", lastError)
}

func (client *Client) CheckApiAccess(in []byte) ([]byte, error) {

	ctx, cancel := context.WithCancel(context.Background())
	result, err := client.grpcClient.CheckApiAccess(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if result.Error != nil {
		cancel()
		return nil, errors.New(string(result.Error))
	}

	cancel()
	return result.Payload, nil

}

func (client *Client) ManageApiAccess(in []byte) ([]byte, error) {

	ctx, cancel := context.WithCancel(context.Background())
	result, err := client.grpcClient.ManageApiAccess(ctx, &pb.Request{Payload: in})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("CLIENT: %v", err)
	}
	if result.Error != nil {
		cancel()
		return nil, errors.New(string(result.Error))
	}

	cancel()
	return result.Payload, nil

}
