//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//
// client.grpcClient.CreateClient
// client.grpcClient.UpdateClient
// client.grpcClient.ListClient
// client.grpcClient.ListClients
// client.grpcClient.RemoveClient

package parallelcore_client_sdk_go

import (
	"context"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

/******************************************************************************/
/*                                                                            */
/* Create Client                                                              */
/*                                                                            */
/******************************************************************************/
func (client *Client) CreateClient(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.CreateClient(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, "CreateClient")
}

/******************************************************************************/
/*                                                                            */
/* Update Client                                                              */
/*                                                                            */
/******************************************************************************/
func (client *Client) UpdateClient(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.UpdateClient(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, "UpdateClient")
}

/******************************************************************************/
/*                                                                            */
/* List Client                                                                */
/*                                                                            */
/******************************************************************************/
func (client *Client) ListClient(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListClient(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, "ListClient")
}

/******************************************************************************/
/*                                                                            */
/* List Clients                                                               */
/*                                                                            */
/******************************************************************************/
func (client *Client) ListClients(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListClients(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, "ListClients")
}

/******************************************************************************/
/*                                                                            */
/* Remove Client                                                              */
/*                                                                            */
/******************************************************************************/
func (client *Client) RemoveClient(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.RemoveClient(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, "RemoveClient")
}
