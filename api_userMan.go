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
/* UserMan                                                                    */
/*                                                                            */
/******************************************************************************/
func (client *Client) userMan(in []byte) ([]byte, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := client.grpcClient.UserMan(ctx, &pb.Request{Payload: in})
	return handleResponse(response, err, API_USER_MAN)
}
