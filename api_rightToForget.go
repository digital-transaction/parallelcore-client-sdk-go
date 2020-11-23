//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//
// client.SysMan(Action: "requestForget"}
// client.SysMan(Action: "approveForget"}
// client.SysMan(Action: "commitForget"}
// client.SysMan(Action: "listForgetGroups"}

package parallelcore_client_sdk_go

import (
	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

/******************************************************************************/
/*                                                                            */
/* Request Forget                                                             */
/*                                                                            */
/******************************************************************************/
// RequestForget request deletion of chunkset(s) covering the given txids.
// It returns a requestTxId which is the transaction id for calling ApproveForget & CommitForget.
func (client *Client) RequestForget(txIds []string) (string, error) {
	x, err := callSysManV(client, API_REQUEST_FORGET, pb.RequestForgetParams{TxIds: txIds})
	return string(x), err
}

/******************************************************************************/
/*                                                                            */
/* Approve Forget                                                             */
/*                                                                            */
/******************************************************************************/
// ApproveForget approve deletion of the above forget request.
// It returns a approvalTxId which provide approvals for calling CommitForget.
func (client *Client) ApproveForget(requestTxId string) (string, error) {
	x, err := callSysManV(client, API_APPROVE_FORGET, pb.ApproveForgetParams{RequestTxId: requestTxId})
	return string(x), err
}

/******************************************************************************/
/*                                                                            */
/* Commit Forget                                                              */
/*                                                                            */
/******************************************************************************/
// CommitForget performs the actual deletion of chunkset based on the
// information from the given request transaction id.
func (client *Client) CommitForget(requestTxId string, approvalTxids []string) (x pb.ForgetReport, err error) {
	_, err = callSysManVV(client, API_COMMIT_FORGET, pb.CommitForgetParams{RequestTxId: requestTxId, ApprovalTxIds: approvalTxids}, &x)
	return x, err
}

/******************************************************************************/
/*                                                                            */
/* List Forget Group                                                          */
/*                                                                            */
/******************************************************************************/
func (client *Client) ListForgetGroups(txIds []string) (x []pb.ForgetGroup, err error) {
	_, err = callSysManVV(client, API_LIST_FORGET_GROUPS, txIds, &x)
	return x, err
}
