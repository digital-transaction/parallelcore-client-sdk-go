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

// certPath: check OpenAny usage
// expireTimestamp: could be set to -1, if user do not want to track when to renew the client

func OpenAnyByToken(endpointSpecs string, token string, expireTimestamp int64, certPath string) (*Client, error) {
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

		client, err := openOne(endpoint, certPath, token)
		if err != nil {
			lastError = err
			endpoints = append(endpoints[:i], endpoints[i+1:]...)
			continue
		}
		client.expireTimestamp = expireTimestamp
		client.endpointSpecs = endpointSpecs

		return client, nil
	}

	return nil, fmt.Errorf("CLIENT: OpenAnyByToken(%q): Failed to open any client using token. Last error: %w", endpointSpecs, lastError)
}

func OpenAnyByTokenWithCert(endpointSpecs string, token string, expireTimestamp int64) (*Client, error) {
	return OpenAnyByToken(endpointSpecs, token, expireTimestamp, os.Getenv("PCORE_CERT_PATH"))
}
