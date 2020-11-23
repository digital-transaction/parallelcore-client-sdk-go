//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"fmt"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"

	"google.golang.org/grpc"
)

type Client struct {
	conn            *grpc.ClientConn
	grpcClient      pb.RequestHandlerClient
	endpointSpecs   string
	certPath        string
	token           string
	expireTimestamp int64
}

func handleResponse(response *pb.Response, err error, function string) ([]byte, error) {
	if err != nil {
		return nil, fmt.Errorf(E_FUNC_X_ERROR_X, function, err)
	}
	if len(response.Error) != 0 {
		return nil, fmt.Errorf(string(response.Error))
	}
	return response.Payload, nil
}

func handleIdentifiedResponse(response *pb.IdentifiedResponse, err error, function string) ([]byte, string, error) {
	if err != nil {
		return nil, "", fmt.Errorf(E_FUNC_X_ERROR_X, function, err)
	}
	if len(response.Error) != 0 {
		return nil, "", fmt.Errorf(string(response.Error))
	}
	return response.Payload, string(response.CommittedId), nil
}
