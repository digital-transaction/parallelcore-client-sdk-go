//
// Copyright 2019 Digital Transaction Limited. All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	"encoding/json"
	"fmt"
)

// ForgetGroup
type ForgetGroup struct {
	TxIds []string `json:"tx_ids"`
}

// ForgetReport
type ForgetReport struct {
	Deleted        []string `json:"deleted"`
	AlreadyDeleted []string `json:"already_deleted"`
	NotFound       []string `json:"not_found"`
	CommitTxId     string   `json:"commit_tx_id"`
}

type requstForgetParams struct {
	TxIds []string `json:"tx_ids"`
}

type approveForgetParams struct {
	RequestTxId string `json:"request_tx_id"`
}

type commitForgetParams struct {
	RequestTxId   string   `json:"request_tx_id"`
	ApprovalTxIds []string `json:"approval_tx_ids"`
}

type forgetTask struct {
	Action string `json:"action"`
	Data   []byte `json:"data"`
}

// RequestForget request deletion of chunkset(s) covering the given txids.
// It returns a requestTxId which is the transaction id for calling ApproveForget & CommitForget.
func (client *Client) RequestForget(txIds []string) (string, error) {
	data, err := json.Marshal(requstForgetParams{TxIds: txIds})
	if err != nil {
		return "", fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	task, err := json.Marshal(forgetTask{Action: "requestForget", Data: data})
	if err != nil {
		return "", fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	requestTxId, err := client.SysMan(task)
	if err != nil {
		return "", err
	}

	return string(requestTxId), nil
}

// ApproveForget approve deletion of the above forget request.
// It returns a approvalTxId which provide approvals for calling CommitForget.
func (client *Client) ApproveForget(requestTxId string) (string, error) {
	data, err := json.Marshal(approveForgetParams{RequestTxId: requestTxId})
	if err != nil {
		return "", fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	task, err := json.Marshal(forgetTask{Action: "approveForget", Data: data})
	if err != nil {
		return "", fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	approvalTxid, err := client.SysMan(task)
	if err != nil {
		return "", err
	}

	return string(approvalTxid), nil
}

// CommitForget performs the actual deletion of chunkset based on the
// information from the given request transaction id.
func (client *Client) CommitForget(requestTxId string, approvalTxids []string) (ForgetReport, error) {
	var report ForgetReport

	data, err := json.Marshal(commitForgetParams{RequestTxId: requestTxId, ApprovalTxIds: approvalTxids})
	if err != nil {
		return report, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	task, err := json.Marshal(forgetTask{Action: "commitForget", Data: data})
	if err != nil {
		return report, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	result, err := client.SysMan(task)
	if err != nil {
		return report, err
	}

	err = json.Unmarshal(result, &report)
	if err != nil {
		return report, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	return report, nil
}

func (client *Client) ListForgetGroups(txIds []string) ([]ForgetGroup, error) {
	var groups []ForgetGroup

	data, err := json.Marshal(txIds)
	if err != nil {
		return groups, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	task, err := json.Marshal(forgetTask{Action: "listForgetGroups", Data: data})
	if err != nil {
		return groups, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	result, err := client.SysMan(task)
	if err != nil {
		return groups, err
	}

	err = json.Unmarshal(result, &groups)
	if err != nil {
		return groups, fmt.Errorf("CLIENT: %s\n", err.Error())
	}

	return groups, nil
}
