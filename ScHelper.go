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

func CallSmartContractJSON(x *Client, name string, action string, v interface{}, result interface{}) (raw []byte, err error) {
	raw, err = CallSmartContract(x, name, action, v)
	if err == nil {
		err = json.Unmarshal(raw, result)
	}
	return
}

func CallSmartContractText(x *Client, name string, action string, v interface{}) (text string, err error) {
	var raw []byte
	raw, err = CallSmartContract(x, name, action, v)
	if err == nil {
		text = string(raw)
	}
	return
}

func CallSmartContract(x *Client, name string, action string, v interface{}) ([]byte, error) {
	var data string
	if v == nil {
		data = ""
	} else if text, ok := v.(string); ok {
		data = text
	} else if raw, ok := v.([]byte); ok {
		data = string(raw)
	} else if raw, e := json.Marshal(v); e == nil {
		data = string(raw)
	} else {
		data = fmt.Sprintf("%v", v)
	}
	task, err := json.Marshal(pb.ScTask{Action: action, Data: data})
	if err != nil {
		return nil, err
	}
	return x.Invoke(name+"-v*", task)
}

func ReturnBytesToString(input []byte, err error) (string, error) {
	return string(input), err
}

func ReturnStringToBytes(input string, err error) ([]byte, error) {
	return []byte(input), err
}
