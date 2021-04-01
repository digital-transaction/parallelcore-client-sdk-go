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

// CreateDomain takes in a string domainName and creates a domain with that name. It returns
// the boolean true if the operation was successful, and (false, error) otherwise.
//
// Permissions: Only super-admins
func (client *Client) CreateDomain(domainName []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.CreateDomain(ctx, &pb.Request{Payload: domainName})
	return handleResponse(response, err, "CreateDomain")
}

// ListDomain (this function could alternatively be named GetDomainInfo) takes in a string domainName (which
// could be empty), and returns information about the specified domain(s).
//
// If domainName is non-empty, ListDomain returns an object with fields:
//  - clients Array of string: Clients that are part of the domain.
//  - admin Array of string: Admins that are part of the domain.
//  - smartContract Array of string: Smart contracts that are part of the domain.
//
// If domainName is empty, ListDomain returns an array of objects with fields:
//  - domainName string
//  - data Object:
//    + clients Array of string: as above.
//    + admin Array of string
//    + smartContract Array of string
//
// ListDomain will only return information about the domains in which the user is admin.
//
// Permissions: Only super-admins and domain-admins
func (client *Client) ListDomain(domainName []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListDomain(ctx, &pb.Request{Payload: domainName})
	return handleResponse(response, err, "ListDomain")
}

// ListManagedDomains takes in a string userID and returns a JSON-encoded array of strings containing
// the names of the domains in which the specified user is admin.
//
// If the user making the call (might not be userID) is a super-admin, they can call
// this function passing in any userID.
//
// If the user making the call is a domain-admin, they need to pass their own userID as
// the argument (non super-admins cannot list the managed domains of other users).
//
// Permissions: Only super-admins and domain-admins
func (client *Client) ListManagedDomains(userID []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListManagedDomains(ctx, &pb.Request{Payload: userID})
	return handleResponse(response, err, API_LIST_MANAGED_DOMAINS)
}

// GrantDomainAdmin takes in a JSON-encoded object with fields:
//  - clientId string
//  - clientDomainName string
//
// and makes the specified client an admin in the specified domain. This is only
// possible if the client is in the target domain in the first place.
//
// Super-admins can GrantDomainAdmin on all domains, domain-admins can only do so
// on domains they manage.
//
// Permissions: Only super-admins and domain-admins
func (client *Client) GrantDomainAdmin(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.GrantDomainAdmin(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_GRANT_DOMAIN_ADMIN)
}

// RevokeDomainAdmin is similar to GrantDomainAdmin, taking in the same parameters and imposing
// the same constraints.
//
// It revokes the specified client's admin privileges in the specified domain.
//
// Permissions: Only super-admins and domain-admins
func (client *Client) RevokeDomainAdmin(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.RevokeDomainAdmin(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_REVOKE_DOMAIN_ADMIN)
}
