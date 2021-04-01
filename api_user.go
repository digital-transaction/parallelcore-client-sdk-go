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
	"encoding/json"
	"strings"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

// CreateClient (deprecated in favor of CreateUser) registers a new
// user in the ParallelChain network. It takes in a JSON-encoded object with fields:
//  - clientId string
//  - clientCredential string
//  - clientRoles string
//  - clientDomainName Array of string or empty string: if empty defaults to 'default' domain
//
// If the calling user is not a super-admin, they can only create clients in domains
// that they manage.
//
// Permissions: Only super-admins and domain-admins
func (client *Client) CreateClient(clientDataJSON []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.CreateClient(ctx, &pb.Request{Payload: clientDataJSON})
	return handleResponse(response, err, "CreateClient")
}

// CreateUser registers a new user in the ParallelChain network. It takes in parameters:
//  - clientID string
//  - password string
//  - roles []string: cannot be empty
//  - domains []string: defaults to ['default'] if empty.
//
// If the calling user is not a super-admin, they can only create clients in domains that
// they manage.
//
// Permissions: Only super-admins and domain-admins.
func (client *Client) CreateUser(userID string, password string, roles []string, domains []string) (string, error) {
	userData := UserData{
		ID:         userID,
		Credential: password,
		Roles:      strings.Join(roles, ","),
		DomainName: strings.Join(domains, ","),
	}
	userDataJSON, err := json.Marshal(userData)

	resRaw, err := client.CreateClient(userDataJSON)
	return string(resRaw), err
}

// UpdateClient (deprecated in favor of UpdateUser) is similar to CreateClient,
// with the same parameters, but updates an existing user instead of creating a new one.
//
// However, it enforces the added constraint that the calling user can only update users
// in domains that they manage (super-admins do not automatically a domain-admin
// of every domain).
//
// Permissions: Only domain-admins
func (client *Client) UpdateClient(clientDataJSON []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.UpdateClient(ctx, &pb.Request{Payload: clientDataJSON})
	return handleResponse(response, err, "UpdateClient")
}

// UpdateUser is similar to CreateUser, with the same parameters, but updates an existing user instead
// of creating a new one.
//
// However, it enforces the added constraint that the calling user can only update
// user in domains that they manage (super-admins do not automatically become a domain-admin
// of every domain).
//
// Permissions: only domain-admins
func (client *Client) UpdateUser(userID string, password string, roles []string, domains []string) (string, error) {
	userData := UserData{
		ID:         userID,
		Credential: password,
		Roles:      strings.Join(roles, ","),
		DomainName: strings.Join(domains, ","),
	}
	userDataJSON, err := json.Marshal(userData)

	resRaw, err := client.UpdateClient(userDataJSON)
	return string(resRaw), err
}

// ListClient (deprecated in favor of GetUserInfo) takes in a string clientID (which could
// be empty), and returns information about the specified user. It returns a JSON-encoded
// object with fields:
//  - clientId string
//  - roles Array of string
//  - accessList Array of string
//  - domains Array of string
//
// If clientID is empty ListClient will return information about the calling user.
//
// Super-admins can ListClient any clientID, domain-admins can ListClient any client in domains
// that they manage, non-admins can ListClient only themselves.
func (client *Client) ListClient(clientID []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListClient(ctx, &pb.Request{Payload: clientID})
	return handleResponse(response, err, "ListClient")
}

// GetUserInfo takes in a string clientID and returns information about the specified user.
// clientID could be left empty, in which case it returns information about the calling user.
// It returns a populated pb.ClientFullData (see type definition).
//
// Super-admins can GetUserInfo any
func (client *Client) GetUserInfo(clientID string) (UserFullData, error) {
	res, err := client.ListClient([]byte(clientID))
	if err != nil {
		return UserFullData{}, err
	}

	resUnmarshaled := UserFullData{}
	if err := json.Unmarshal(res, &resUnmarshaled); err != nil {
		return UserFullData{}, err
	}

	return resUnmarshaled, err
}

// ListClients (deprecated in favor of GetUserInfos) takes in a JSON-encoded object
// with keys:
//  - allDomains bool
//  - domainName string
//
// It returns a JSON-encoded array of objects with fields:
//  - clientID string
//  - data Object:
//    + clientID string
//    + roles Array of string
//    + accessList Array of string
//    + domains Array of string
//    + managedDomains Array of string
//
// Domain-admins can only list users in domains that they manage.
//
// Permissions: Only super-admins or domain-admins
func (client *Client) ListClients(query []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListClients(ctx, &pb.Request{Payload: query})
	return handleResponse(response, err, "ListClients")
}

// GetUserInfos is similar to GetUserInfo, but returns []UserFullDataWrapper (see type definition),
// containing information about multiple users. It takes in two parameters:
//  - allDomains bool
//  - domainName string
//
// There are two intended uses for GetUserInfos:
//  - allDomains == true && domainName == "": returns information about all users.
//  - allDomains == false && domainName == <domainName>: returns information about users in specified domain.
//
// Mere domain-admins can only GetUserInfos in domains that they manage.
//
// Permissions: Only super-admins or domain-admins
func (client *Client) GetUserInfos(allDomains bool, domainName string) ([]UserFullDataWrapper, error) {
	query, _ := json.Marshal(InfoListData{
		AllDomains: allDomains,
		DomainName: domainName,
	})

	resRaw, err := client.ListClients([]byte(query))
	if err != nil {
		return []UserFullDataWrapper{}, err
	}

	resUnmarshaled := []UserFullDataWrapper{}
	if err := json.Unmarshal(resRaw, &resUnmarshaled); err != nil {
		return []UserFullDataWrapper{}, err
	}

	return resUnmarshaled, err
}

// RemoveClient (deprecated in favor of DeleteUser) takes in a JSON-encoded object
// with keys:
//  - clientId string
//  - clientDomainName string
//
// and completely removes the user identified by clientId from the network. The calling user
// can only remove users in domains that they manage, and cannot remove themselves. clientDomainName
// can be the name of any domain that the user is part of, and defaults to 'default' if it is
// an empty string.
//
// Permissions: Only super-admins or domain-admins
//
// BUG(RemoveClient): RemoveClient: non super-admins should not be allowed to remove a user from the network
// entirely. Presently, 'mere' domain-admins are allowed to do this.
func (client *Client) RemoveClient(clientDomainDataJSON []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.RemoveClient(ctx, &pb.Request{Payload: clientDomainDataJSON})
	return handleResponse(response, err, "RemoveClient")
}

// DeleteUser completely removes the user identified by userID from the network. The calling user
// can only remove users in domains that they manage, and cannot remove themselves.  userDomainName
// can be the name of any domain that the user is part of, and defaults to 'default' if it is an
// empty string.
//
// Permissions: Only super-admins or domain-admins
//
// BUG(DeleteUser): DeleteUser: non super-admins should not be allowed to remove a user from the network
// entirely. Presently, 'mere' domain-admins are allowed to do this.
func (client *Client) DeleteUser(userID string, userDomainName string) (string, error) {
	userDomainData := UserDomainData{
		ID:         userID,
		DomainName: userDomainName,
	}
	userDomainDataJSON, _ := json.Marshal(userDomainData)

	resRaw, err := client.RemoveClient(userDomainDataJSON)
	return string(resRaw), err
}
