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

// GrantAccess is used to grant users access to (i.e. have the permission to invoke) smart contracts.
// It takes in a JSON-encoded object with fields:
//  - clientId string
//  - scName string
//  - domainName string: a domain that the client is a part of
//
// The calling user can only GrantAccess to smart contracts registered in domains that they manage, to
// users in domains that they manage. Note that in ParallelChain, there is no concept of 'invoking an SC
// in a domain.' Invoking an SC is always invoking an SC simpliciter (apologies for the philosophical
// language).
//
// Permissions: Only super-admins and domain-admins
func (client *Client) GrantAccess(clientAccessDataJSON []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.GrantAccess(ctx, &pb.Request{Payload: clientAccessDataJSON})
	return handleResponse(response, err, API_GRANT_ACCESS)
}

// TODO: Adrio
// func (client *Client) GrantAccessToSC(clientID string, scName string, domainName string)

// RevokeAccess is similar to GrantAccess, with the same parameters, but revokes access to the specified
// scName instead of granting access.
//
// Permissions: Only super-admins and domain-admins
func (client *Client) RevokeAccess(clientAccessDataJSON []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.RevokeAccess(ctx, &pb.Request{Payload: clientAccessDataJSON})
	return handleResponse(response, err, API_REVOKE_ACCESS)
}

// TODO: Adrio
// func (client *Client) RevokeAccessToSC(clientID string, scName string, domainName string)

// CheckApiAccess is used to check if the calling user has access to (can use) a specified functionality
// provided by the ParallelCore RESTful API. It takes in a JSON-encoded object with fields:
//  - operation string
//  - api string
//  - options ???
//
// BUG(CheckApiAccess): CheckApiAccess: we do not know what this does exactly.
func (client *Client) CheckApiAccess(apiAccessControllerJSON []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.CheckApiAccess(ctx, &pb.Request{Payload: apiAccessControllerJSON})
	return handleResponse(response, err, API_CHECK_API_ACCESS)
}

// ManageApiAccess is similar to CheckApiAccess.
func (client *Client) ManageApiAccess(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ManageApiAccess(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_MANAGE_API_ACCESS)
}
