//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"context"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

/******************************************************************************/
/*                                                                            */
/* Register Smart Contract                                                    */
/*                                                                            */
/******************************************************************************/
func (client *Client) RegisterSmartContract(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.RegisterSmartContract(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_REGISTER_SMARTCONTRACT)
}

/******************************************************************************/
/*                                                                            */
/* List Smart Contract                                                        */
/*                                                                            */
/******************************************************************************/
func (client *Client) ListSmartContract(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListSmartContract(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_LIST_SMARTCONTRACT)
}

/******************************************************************************/
/*                                                                            */
/* List Smart Contracts                                                       */
/*                                                                            */
/******************************************************************************/
func (client *Client) ListSmartContracts(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.ListSmartContracts(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_LIST_SMARTCONTRACTS)
}
