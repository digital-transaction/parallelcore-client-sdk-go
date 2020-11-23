//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//
// client.grpcClient.GrantAccess
// client.grpcClient.RevokeAccess
// client.grpcClient.CheckApiAccess
// client.grpcClient.ManageApiAccess

package parallelcore_client_sdk_go

import (
	"context"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

/******************************************************************************/
/*                                                                            */
/* Grant Access                                                               */
/*                                                                            */
/******************************************************************************/
func (client *Client) GrantAccess(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.GrantAccess(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_GRANT_ACCESS)
}

/******************************************************************************/
/*                                                                            */
/* Revoke Access                                                              */
/*                                                                            */
/******************************************************************************/
func (client *Client) RevokeAccess(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.RevokeAccess(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_REVOKE_ACCESS)
}

/******************************************************************************/
/*                                                                            */
/* Check API Access                                                           */
/*                                                                            */
/******************************************************************************/
func (client *Client) CheckApiAccess(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.CheckApiAccess(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_CHECK_API_ACCESS)
}

/******************************************************************************/
/*                                                                            */
/* Manage API Access                                                          */
/*                                                                            */
/******************************************************************************/
func (client *Client) ManageApiAccess(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ManageApiAccess(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_MANAGE_API_ACCESS)
}
