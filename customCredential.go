//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"context"
)

// customCredential, which stores the JWT
type customCredential struct {
	token string
}

func (t customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (t customCredential) RequireTransportSecurity() bool {
	return true
}
