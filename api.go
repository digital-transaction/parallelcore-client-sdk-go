// Package parallelcore_client_sdk_go (The ParallelCore Go Client SDK)
// enables programmers to develop applications that interact with
// ParallelChain networks through invocation of smart contracts.
//
// Basic workflow:
//  1. Use OpenAny() to establish a connection from the application to ParallelCore gRPC endpoint(s).
//  2. Invoke a smart-contract using Invoke(), passing in arguments as a space-delimited string.
//  3. After the application finishes using the connection, close it using Close()
//
// Copyright 2021 Digital Transaction Limited.
// All Rights Reserved.
package parallelcore_client_sdk_go

import (
	"fmt"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"

	"google.golang.org/grpc"
)

// Client represents a connection to a ParallelCore node.
//
// It provides all of the methods an credentialed 'application' (program running
// outside of the blockchain network; for example in a desktop, website,
// or mobile phone) can use to invoke smart contracts, perform system
// administration, inspect the blockchain, etc.
//
// The typical procedure to instantiate a Client is to call
// ClientSDK.OpenAny or ClientSDK.OpenMany.
type Client struct {
	conn            *grpc.ClientConn
	grpcClient      pb.RequestHandlerClient
	endpointSpecs   string
	certPath        string
	token           string
	expireTimestamp int64
}

/*
func (x *Client) String() string {
	return fmt.Sprintf("PCoreClient(token:'%s...%s' expire:%s)", x.token[:4], x.token[len(x.token)-4:], time.Unix(x.expireTimestamp, 0).Format("2006/01/02_15:04:05"))
}
*/

var endpointsCalled = 0
var oldEndpoints string

func handleResponse(response *pb.Response, err error, function string) ([]byte, error) {
	if err != nil {
		return nil, fmt.Errorf(E_FUNC_X_ERROR_X, function, err)
	}
	if len(response.Error) != 0 {
		return nil, fmt.Errorf(string(response.Error))
	}
	return response.Payload, nil
}

func handleIdentifiedResponse(response *pb.IdentifiedResponse, err error, function string) ([]byte, string, error) {
	if err != nil {
		return nil, "", fmt.Errorf(E_FUNC_X_ERROR_X, function, err)
	}
	if len(response.Error) != 0 {
		return nil, "", fmt.Errorf(string(response.Error))
	}
	return response.Payload, string(response.CommittedId), nil
}
