//
// Copyright DTL. All Rights Reserved.
//
// Authors: Yang SONG, Eric Ma
//

package engine_client_proto

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"net"
	"strings"
)

// Utils used by servers for client-server interaction

func StartServer(lis net.Listener, srv RequestHandlerServer, enableTLS bool, pemPath string, keyPath string) error {

	if enableTLS {

		// TLS
		fmt.Println("With TLS.")
		creds, err := credentials.NewServerTLSFromFile(pemPath, keyPath)
		if err != nil {
			return fmt.Errorf("Failed to generate credentials. Error:  %v", err)
		}

		grpcServer := grpc.NewServer(grpc.Creds(creds))
		RegisterRequestHandlerServer(grpcServer, srv)
		fmt.Println("Running!")
		grpcServer.Serve(lis)

	} else {

		// Without TLS
		fmt.Println("Without TLS.")
		grpcServer := grpc.NewServer()
		RegisterRequestHandlerServer(grpcServer, srv)
		fmt.Println("Running!")
		grpcServer.Serve(lis)

	}

	return nil
}

func ParseContext(ctx context.Context) (string, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("Context parsing error.")
	}

	var token string
	if val, ok := md["authorization"]; ok {
		token = val[0]
	}
	if token == "" {
		return "", fmt.Errorf("Token not found.")
	}

	tmp := strings.Split(token, "Bearer ")
	if len(tmp) != 2 {
		return "", fmt.Errorf("Error while parsing token")
	}
	token = tmp[1]

	return token, nil
}
