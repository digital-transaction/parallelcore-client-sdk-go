//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//
// client.grpcClient.Auth

package parallelcore_client_sdk_go

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

func (client *Client) auth(clientID []byte, credential []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.Auth(ctx, &pb.AuthRequest{ClientId: clientID, Credential: credential})
	return handleResponse(response, err, "auth")
}

// UpdateSelfCredential is used to change the credential (password) of the user identified by
// clientId. Users (even super-admins) can only change their own credentials. credential cannot
// be an empty string.
func (client *Client) UpdateSelfCredential(clientID string, credential string) ([]byte, error) {
	return callUserManV(client, API_UPDATE_SELF_CREDENTIAL, UserData{ID: clientID, Credential: credential})
}

// GetToken is a getter for the private field client.token.
//
// client.token is used for token-based authentication using the Open*ByToken group of
// functions.
func (client *Client) GetToken() string {
	return client.token
}

// GetTokenExpTime is a getter for the private field client.expireTimestamp.
//
// client.expireTimestamp is used by applications to decide when to call Renew().
func (client *Client) GetTokenExpTime() int64 {
	return client.expireTimestamp
}

// Renew asks a ParallelCore endpoint randomly selected from client.endpointSpecs
// to renew the calling client's authentication token.
func (client *Client) Renew() error {
	token, expireTimestamp, err := client.renewToken()
	if err != nil {
		return err
	}

	client.Close()

	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	endpoints := strings.Split(client.endpointSpecs, " ")
	var lastError error
	for len(endpoints) != 0 {
		i := randGen.Intn(len(endpoints))
		endpoint := endpoints[i]

		newClient, err := openOne(endpoint, client.certPath, token)
		if err != nil {
			lastError = err
			endpoints = append(endpoints[:i], endpoints[i+1:]...)
			continue
		}
		newClient.expireTimestamp = expireTimestamp
		newClient.endpointSpecs = client.endpointSpecs

		*client = *newClient
		return nil
	}
	return fmt.Errorf("CLIENT: Failed to renew a connection. %v", lastError)
}

func (client *Client) renewToken() (token string, expireTimestamp int64, err error) {
	var response *pb.Response
	// Fetch new JWT and expireTimestamp
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err = client.grpcClient.Renew(ctx, &pb.Request{Payload: []byte("")})
	if err != nil {
		return "", 0, fmt.Errorf(E_FUNC_X_ERROR_X, "renewToken", err)
	}
	if len(response.Error) != 0 {
		return "", 0, fmt.Errorf(FMT_FUNC_X_RESPONSE_ERROR_X, "renewToken", fmt.Errorf(string(response.Error)))
	}
	// Parse returnBytes
	token, expireTimestamp, err = parseTokenAndExpireTimestamp(string(response.Payload))
	if err != nil {
		err = fmt.Errorf("CLIENT: %w", err)
	}
	return token, expireTimestamp, err
}
