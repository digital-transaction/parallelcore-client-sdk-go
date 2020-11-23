//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"fmt"
	"os"
	"strings"
)

// certPath: check OpenAny usage

func OpenMany(endpointSpecs string, clientId string, credential string, certPath string) ([]*Client, error) {
	endpoints := strings.Split(endpointSpecs, " ")
	lastError := fmt.Errorf("")
	clients := make([]*Client, 0)
	// these clients can use the same token
	token := ""
	var expireTimestamp int64
	for _, endpoint := range endpoints {
		// Create a connection first.
		// If the connection is not connect with token, then will run next block to fetch JWT and expireTimestamp
		client, err := openOne(endpoint, certPath, token)
		if err != nil {
			lastError = fmt.Errorf("CLIENT: OpenMany(%q): %w. Last error: %v", endpointSpecs, err, lastError)
			continue
		}
		if token == "" {
			// Successfully setup the first connection, fetch the token
			returnBytes, err := client.auth([]byte(clientId), []byte(credential))
			client.Close()
			if err != nil {
				return nil, fmt.Errorf("CLIENT: OpenMany(%q): Failed to auth. %w", endpointSpecs, err)
			}

			// Parse returnBytes
			token, _, err := parseTokenAndExpireTimestamp(string(returnBytes))
			if err != nil {
				return nil, fmt.Errorf("CLIENT: OpenMany(%q): Failed to parse Token. %w", endpointSpecs, err)
			}

			client, err = openOne(endpoint, certPath, token)
			if err != nil {
				return nil, fmt.Errorf("CLIENT: OpenMany(%q): Failed to open a connection. %w", endpointSpecs, err)
			}
		}
		client.expireTimestamp = expireTimestamp
		client.endpointSpecs = endpointSpecs
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		// All attemps to connect are failed.
		return nil, fmt.Errorf("CLIENT: OpenMany(%q): Failed to open all clients. Last error: %w", endpointSpecs, lastError)
	}
	if lastError.Error() == "" {
		return clients, nil
	}
	// All good or some attemps to connect failed. Opened all clients or opened some clients.
	return clients, lastError
}

func OpenManyWithCert(endpointSpecs string, clientId string, credential string) ([]*Client, error) {
	return OpenMany(endpointSpecs, clientId, credential, os.Getenv("PCORE_CERT_PATH"))
}

func CloseMany(clients []*Client) {
	for _, each := range clients {
		each.Close()
	}
}
