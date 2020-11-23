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

func OpenManyByToken(endpointSpecs string, token string, expireTimestamp int64, certPath string) ([]*Client, error) {
	endpoints := strings.Split(endpointSpecs, " ")
	lastError := fmt.Errorf("")
	clients := make([]*Client, 0)
	for _, endpoint := range endpoints {
		client, err := openOne(endpoint, certPath, token)
		if err != nil {
			lastError = fmt.Errorf("CLIENT: OpenManyByToken(%q): %w. Last error: %v", endpointSpecs, err, lastError)
			continue
		}
		client.expireTimestamp = expireTimestamp
		client.endpointSpecs = endpointSpecs
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		// All attemps to connect are failed.
		return nil, fmt.Errorf("CLIENT: OpenManyByToken(%q): Failed to open all clients using token. Last error: %w", endpointSpecs, lastError)
	}
	if lastError.Error() == "" {
		return clients, nil
	}
	// All good or some attemps to connect failed. Opened all clients or opened some clients.
	return clients, lastError
}

func OpenManyByTokenWithCert(endpointSpecs string, token string, expireTimestamp int64) ([]*Client, error) {
	return OpenManyByToken(endpointSpecs, token, expireTimestamp, os.Getenv("PCORE_CERT_PATH"))
}
