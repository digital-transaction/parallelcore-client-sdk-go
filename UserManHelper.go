//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"encoding/json"
	"fmt"

	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

func callUserMan(client *Client, action string, data []byte) ([]byte, error) {
	// Encode Task
	task, err := json.Marshal(pb.UserManData{Action: action, Data: data})
	if err != nil {
		return nil, fmt.Errorf(FMT_FUNC_X_TASK_ENCODE_ERROR_X, action, err)
	}
	// Call Task
	return client.userMan(task)
}

func callUserManV(client *Client, action string, v interface{}) ([]byte, error) {
	// Encode Data
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf(FMT_FUNC_X_INPUT_ENCODE_ERROR_X, action, err)
	}
	// Call Task
	return callUserMan(client, action, data)
}
