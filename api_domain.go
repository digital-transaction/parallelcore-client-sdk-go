//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//
// client.grpcClient.CreateDomain
// client.grpcClient.ListDomain
// client.grpcClient.ListManagedDomains
// client.grpcClient.GrantDomainAdmin
// client.grpcClient.RevokeDomainAdmin

package parallelcore_client_sdk_go

import (
	"context"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

/******************************************************************************/
/*                                                                            */
/* Create Domain                                                              */
/*                                                                            */
/******************************************************************************/
func (client *Client) CreateDomain(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.CreateDomain(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, "CreateDomain")
}

/******************************************************************************/
/*                                                                            */
/* List Domain                                                                */
/*                                                                            */
/******************************************************************************/
func (client *Client) ListDomain(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListDomain(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, "ListDomain")
}

/******************************************************************************/
/*                                                                            */
/* List Manage Domains                                                        */
/*                                                                            */
/******************************************************************************/
func (client *Client) ListManagedDomains(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListManagedDomains(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_LIST_MANAGED_DOMAINS)
}

/******************************************************************************/
/*                                                                            */
/* Grant Domain Admin                                                         */
/*                                                                            */
/******************************************************************************/
func (client *Client) GrantDomainAdmin(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.GrantDomainAdmin(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_GRANT_DOMAIN_ADMIN)
}

/******************************************************************************/
/*                                                                            */
/* Revoke Domain Admin                                                        */
/*                                                                            */
/******************************************************************************/
func (client *Client) RevokeDomainAdmin(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.RevokeDomainAdmin(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_REVOKE_DOMAIN_ADMIN)
}
