//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"context"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

// SysMan (deprecated) is used to make non super-admin level queries about, or changes to, the configuration of
// a ParallelChain network. This includes updating one's password, querying blockchain/transaction
// summary data, and so on.
//
// All of these functionalities are now implemented in dedicated methods.
func (client *Client) userMan(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.UserMan(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_USER_MAN)
}
