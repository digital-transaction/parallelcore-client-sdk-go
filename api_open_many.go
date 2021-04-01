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

// OpenMany is similar to OpenAny, taking in the same parameters, but tries connect to all endpoints instead
// of one. Accordingly, it returns a list of Clients instead of just one. If one connection attempt fails,
// OpenMany skips the endpoint.
//
// If all connection attempts fail, it will return (nil, error). Otherwise, error is the last error
// encountered during connection attempts. Applications should check whether or not []*Client is nil
// to determine if OpenMany succeeded to establish some connections. Do not use error for this purpose.
func OpenMany(endpointSpecs string, clientID string, credential string, certPath string) ([]*Client, error) {
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
			returnBytes, err := client.auth([]byte(clientID), []byte(credential))
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

// OpenManyWithCert is a wrapper around OpenMany. It calls OpenMany with os.Getenv("PCORE_CERT_PATH")
// as the certPath parameter.
func OpenManyWithCert(endpointSpecs string, clientID string, credential string) ([]*Client, error) {
	return OpenMany(endpointSpecs, clientID, credential, os.Getenv("PCORE_CERT_PATH"))
}

// CloseMany is like Close, but closes every Client in clients.
func CloseMany(clients []*Client) {
	for _, each := range clients {
		each.Close()
	}
}
