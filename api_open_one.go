//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//
// grpcClient.Ping

package parallelcore_client_sdk_go

import (
	"crypto/tls"
	"fmt"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

func openOne(endpoint string, certPath string, token string) (_ *Client, err error) {
	var (
		creds credentials.TransportCredentials
		conn  *grpc.ClientConn
	)
	if certPath != "" {
		creds, err = credentials.NewClientTLSFromFile(certPath, "")
		if err != nil {
			return nil, fmt.Errorf("CLIENT: openOne(%q): %w", endpoint, err)
		}
	} else {
		creds = credentials.NewTLS(&tls.Config{})
	}

	grpcOpts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	if token != "" {
		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(customCredential{token: token}))
	}
	// WithBlock returns a DialOption which makes caller of Dial blocks until the underlying connection is up.
	// Without this, Dial returns immediately and connecting the server happens in background.
	grpcOpts = append(grpcOpts, grpc.WithBlock())

	conn, err = grpc.Dial(endpoint, grpcOpts...)
	if err != nil {
		return nil, fmt.Errorf("CLIENT: openOne(%q): Failed to dail. %w", endpoint, err)
	}
	if conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf("CLIENT: openOne(%q): Connection Failed", endpoint)
	}

	grpcClient := pb.NewRequestHandlerClient(conn)

	return &Client{conn, grpcClient, "", certPath, token, 0}, nil
}

func (client *Client) Close() {
	if client.conn != nil {
		client.conn.Close()
	}
}
