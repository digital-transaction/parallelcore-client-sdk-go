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

// certPath: if empty, use the system certificate, otherwise, use the certificate provided in the file in certPath

// OpenAny establishes and returns one client connection to a ParallelChain peer. It takes in:
//  - endpointSpecs string: space-delimited list of endpoints. This may contain however many endpoints,
//  but OpenAny will only choose one to connect to. If a connection attempt fails, it will choose another
//  endpoint.
//  - clientID string
//  - clientCredential string
//  - certPath string: file path to a TLS certificate to set up an encrypted connection. If certPath is empty,
//  the system certificate pool will be used.
func OpenAny(endpointSpecs string, clientID string, credential string, certPath string) (*Client, error) {
	// do until endpoints is empty:
	//  select first endpoint from endpoints
	//  try to openOne(endpoint) -> c
	//  if ok
	//    return c
	//  remove endpoint from endpoints
	// return error
	var lastError error

	if (oldEndpoints == "") || (oldEndpoints != endpointSpecs) {
		oldEndpoints = endpointSpecs
		endpointsCalled = 0
	}

	endpoints := strings.Split(endpointSpecs, " ")
	for len(endpoints) != 0 {
		i := endpointsCalled % len(endpoints)
		endpoint := endpoints[i]

		client, err := openOne(endpoint, certPath, "")
		if err != nil {
			lastError = err
			endpoints = append(endpoints[:i], endpoints[i+1:]...)
			continue
		}
		// Successfully setup a connection, fetch the token and expireTimestamp
		returnBytes, err := client.auth([]byte(clientID), []byte(credential))
		client.Close()
		if err != nil {
			return nil, fmt.Errorf("CLIENT: OpenAny(%q): Failed to auth. %w", endpointSpecs, err)
		}
		// Parse returnBytes
		token, expireTimestamp, err := parseTokenAndExpireTimestamp(string(returnBytes))
		if err != nil {
			return nil, fmt.Errorf("CLIENT: OpenAny(%q): Failed to parse Token. %w", endpointSpecs, err)
		}
		client, err = openOne(endpoint, certPath, token)
		if err != nil {
			return nil, fmt.Errorf("CLIENT: OpenAny(%q): Failed to open a connection. %w", endpointSpecs, err)
		}
		endpointsCalled++
		client.expireTimestamp = expireTimestamp
		client.endpointSpecs = endpointSpecs

		return client, nil
	}

	return nil, fmt.Errorf("CLIENT: OpenAny(%q): Failed to open any client. Last error: %w", endpointSpecs, lastError)
}

// OpenAnyWithCert is a wrapper around OpenAny. It calls OpenAny with os.Getenv("PCORE_CERT_PATH")
// as the certPath parameter.
func OpenAnyWithCert(endpointSpecs string, clientID string, credential string) (*Client, error) {
	return OpenAny(endpointSpecs, clientID, credential, os.Getenv("PCORE_CERT_PATH"))
}
