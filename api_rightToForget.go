//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//
// client.SysMan(Action: "requestForget"}
// client.SysMan(Action: "approveForget"}
// client.SysMan(Action: "commitForget"}
// client.SysMan(Action: "listForgetGroups"}

package parallelcore_client_sdk_go

// RequestForget requests deletion of all transactions specified in txIds.
// The transactions specified in txIds form a 'forget group': either they are all deleted together,
// or they are not deleted.
//
// It returns a string, forgetRequestTxID, which can be used to approve the forget
// request using ApproveForget.
//
// Permissions: Only super-admins.
func (client *Client) RequestForget(txIds []string) (string, error) {
	x, err := callSysManV(client, API_REQUEST_FORGET, RequestForgetParams{TxIds: txIds})
	return string(x), err
}

// ApproveForget approves deletion of a forget request created using RequestFormet
// It returns a string, forgetApprovalTxID, which can be used to commit the forget request
// using CommitForget.
//
// Permissions: Only super-admins.
func (client *Client) ApproveForget(forgetRequestTxID string) (string, error) {
	x, err := callSysManV(client, API_APPROVE_FORGET, ApproveForgetParams{RequestTxId: forgetRequestTxID})
	return string(x), err
}

// CommitForget commits (actually performs) the deletion of all transactions specified
// in the txIds originally passed to the RequestForget call that produced forgetRequestTxID.
//
// It returns a JSON-encoded object with keys:
//  - deleted list: txIds in the forget group actually deleted by the commit
//  - already_deleted list: txIds that were deleted before this commit
//  - not_found: txIds in the forget group that were invalid or not found
//  - commit_tx_id: transaction ID of the commit forget transaction.
//
// Permissions: Only super-admins.
func (client *Client) CommitForget(forgetRequestTxID string, forgetApprovalTxID []string) (x ForgetReport, err error) {
	_, err = callSysManVV(client, API_COMMIT_FORGET, CommitForgetParams{RequestTxId: forgetRequestTxID, ApprovalTxIds: forgetApprovalTxID}, &x)
	return x, err
}

// ListForgetGroups returns a list of all forget groups that cover all txIds (a super-set
// of txIds).
//
// Permissions: Only super-admins.
func (client *Client) ListForgetGroups(txIds []string) (x []ForgetGroup, err error) {
	_, err = callSysManVV(client, API_LIST_FORGET_GROUPS, txIds, &x)
	return x, err
}
