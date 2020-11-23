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

/******************************************************************************/
/*                                                                            */
/* List Invokable Smartcontractgs                                             */
/*                                                                            */
/******************************************************************************/
func (client *Client) ListInvokableSC() ([]byte, error) {
	return callUserMan(client, API_LIST_INVOKABLE_SC, make([]byte, 0))
}

/******************************************************************************/
/*                                                                            */
/* Invoke Smartcontract                                                       */
/*                                                                            */
/******************************************************************************/
func (client *Client) Invoke(smartcontract_spec string, args []byte) ([]byte, error) {
	return client.invoke(append([]byte(smartcontract_spec+" "), args...))
}

/******************************************************************************/
/*                                                                            */
/* Identified Invoke Smartcontract                                            */
/*                                                                            */
/******************************************************************************/
func (client *Client) IdentifiedInvoke(smartcontract_spec string, args []byte) ([]byte, string, error) {
	return client.identifiedInvoke(append([]byte(smartcontract_spec+" "), args...))
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
