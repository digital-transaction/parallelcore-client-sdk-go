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

func callSysMan(client *Client, action string, data []byte) ([]byte, error) {
	// Encode Task
	task, err := json.Marshal(pb.SysManData{Action: action, Data: data})
	if err != nil {
		return nil, fmt.Errorf(FMT_FUNC_X_TASK_ENCODE_ERROR_X, action, err)
	}
	// Call Task
	return client.SysMan(task)
}

func callSysManV(client *Client, action string, v interface{}) ([]byte, error) {
	// Encode Data
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf(FMT_FUNC_X_INPUT_ENCODE_ERROR_X, action, err)
	}
	// Call Task
	return callSysMan(client, action, data)
}

func callSysManVV(client *Client, action string, v interface{}, result interface{}) ([]byte, error) {
	bytesReturn, err := callSysManV(client, action, v)
	if err == nil {
		err = json.Unmarshal(bytesReturn, &result)
	}
	return bytesReturn, err
}
