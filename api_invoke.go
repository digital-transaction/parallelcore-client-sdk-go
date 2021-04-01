//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//
// client.grpcClient.Invoke
// client.grpcClient.IdentifiedInvoke

package parallelcore_client_sdk_go

import (
	"context"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

// ListInvokableSC returns a JSON-encoded array listing the smart contracts
// invokable by the client.
//
// Each item in the array is an object with fields:
//  - name string
//  - ver string
func (client *Client) ListInvokableSC() ([]byte, error) {
	return callUserMan(client, API_LIST_INVOKABLE_SC, make([]byte, 0))
}

// Invoke invokes the smart contract identified by smartContractSpec, blocking until the
// smart contract invocation completes. It returns the byte slice returned by the smart
// contract's Handle method.
//
// It takes in:
//  - smartContractSpec string: SC identifier with the format: <SC name>-v<SC version number>
//  - args string: passed into the invoke SC's Handle function as its 2nd 'in' parameter.
func (client *Client) Invoke(smartContractSpec string, args []byte) ([]byte, error) {
	return client.invoke(append([]byte(smartContractSpec+" "), args...))
}

// IdentifiedInvoke is similar to Invoke but has as its 2nd returned value the SC's transaction
// commit ID. If the transaction the SC produces is a read-only transaction, or if the invocation
// errors, commit ID will be an empty string.
func (client *Client) IdentifiedInvoke(smartContractSpec string, args []byte) ([]byte, string, error) {
	return client.identifiedInvoke(append([]byte(smartContractSpec+" "), args...))
}

func (client *Client) invoke(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.Invoke(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, "invoke")
}

func (client *Client) identifiedInvoke(in []byte) ([]byte, string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.IdentifiedInvoke(ctx, &pb.Request{Payload: in})
	return handleIdentifiedResponse(response, err, "identifiedInvoke")
}
