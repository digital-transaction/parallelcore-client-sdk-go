//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// certPath: if empty, use the system certificate, otherwise, use the certificate provided in the file in certPath

func OpenAny(endpointSpecs string, clientId string, credential string, certPath string) (*Client, error) {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	// do until endpoints is empty:
	//  randomly select endpoint from endpoints
	//  try to openOne(endpoint) -> c
	//  if ok
	//    return c
	//  remove endpoint from endpoints
	// return error

	var lastError error
	endpoints := strings.Split(endpointSpecs, " ")
	for len(endpoints) != 0 {
		i := randGen.Intn(len(endpoints))
		endpoint := endpoints[i]

		client, err := openOne(endpoint, certPath, "")
		if err != nil {
			lastError = err
			endpoints = append(endpoints[:i], endpoints[i+1:]...)
			continue
		}
		// Successfully setup a connection, fetch the token and expireTimestamp
		returnBytes, err := client.auth([]byte(clientId), []byte(credential))
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
		client.expireTimestamp = expireTimestamp
		client.endpointSpecs = endpointSpecs

		return client, nil
	}

	return nil, fmt.Errorf("CLIENT: OpenAny(%q): Failed to open any client. Last error: %w", endpointSpecs, lastError)
}

func OpenAnyWithCert(endpointSpecs string, clientId string, credential string) (*Client, error) {
	return OpenAny(endpointSpecs, clientId, credential, os.Getenv("PCORE_CERT_PATH"))
}
