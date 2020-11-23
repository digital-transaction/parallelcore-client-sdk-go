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

/******************************************************************************/
/*                                                                            */
/* SysMan                                                                     */
/*                                                                            */
/******************************************************************************/
func (client *Client) SysMan(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.SysMan(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_SYS_MAN)
}
