//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"context"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

// RegisterSmartContract is used to install and register smart contracts for invocation
// in the network, whether new or an updated version.
//
// It takes in JSON-encoded object with keys:
//  - scName string: In the form: <scName-scVersion>
//  - FileContent string: Raw contents of Go build-able SC source code packaged into a .tgz
//  - DomainName string: The domain onto which the smart contract will be installed.
//  - InitArgs string: Passed into Initialize method of SC as its second argument.
//
// Users can only register SCs into domains in which they are admin.
//
// Permissions: only super-admins or domain-admins
func (client *Client) RegisterSmartContract(scRegistration []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.RegisterSmartContract(ctx, &pb.Request{Payload: scRegistration})
	return handleResponse(response, err, API_REGISTER_SMARTCONTRACT)
}

// ListSmartContract (this function could alternatively be named GetSmartContractInfo) takes in the name of
// a smart contract and returns information about the SC in a JSON-encoded object with fields:
//  - scName string
//  - scVersion string
//  - space string: namespace (also known as keyspace)
//  - checksum string
//  - mode string: (deprecated) either "binary" or "package"
//  - accessList string: comma-delimited list of clients that can invoke this SC
//  - domains string
//
// Users can only get the info of SCs in domains in which they are admin.
//
// Permissions: Only super-admins or domain-admins
func (client *Client) ListSmartContract(scName []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListSmartContract(ctx, &pb.Request{Payload: scName})
	return handleResponse(response, err, API_LIST_SMARTCONTRACT)
}

// ListSmartContracts takes in a JSON-encoded object with keys:
//  - allDomains bool
//  - domainName string
//
// It returns a JSON-encoded array of objects with information about smart contracts (see
// ListSmartContract (singular) for fields returned).
//
// If allDomains is true, then the array contains information for all SCs in the system (if user
// is super-admin) or all SCs which the user is domain-admin over (if user is not domain-admin)
//
// If allDomains is false, the array contains all smart contracts in domainName.
//
// Permissions: Only super-admins or domain-admins
func (client *Client) ListSmartContracts(query []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListSmartContracts(ctx, &pb.Request{Payload: query})
	return handleResponse(response, err, API_LIST_SMARTCONTRACTS)
}
