//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//
// client.grpcClient.SysMan

package parallelcore_client_sdk_go

import (
	"context"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

// SysMan (deprecated) is used to make super-admin level queries about, or changes to, the configuration of
// a ParallelChain network. This includes creating, updating, removing, and listing clients,
// granting access to smart contracts, requesting, approving, and commiting forget of
// transactions, and so on.
//
// All of these functionalities are now implemented in dedicated methods.
func (client *Client) SysMan(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.SysMan(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_SYS_MAN)
}
