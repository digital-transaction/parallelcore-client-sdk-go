//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

// GetSmartContractTransactionJson returns a JSON-encoded object containing
// the DB-related (list of mutations) details of the transaction identified by transactionID.
// The 'key' and 'value' of a transaction mutation are returned in base64 encoding.
//
// As an example, key: base64("QWxpY2U") === "Alice" in human-readable form.
func (client *Client) GetSmartContractTransactionJson(transactionId string) ([]byte, error) {
	return callUserMan(client, API_GET_SMARTCONTRACT_TRANSACTION_JSON, []byte(transactionId))
}

// GetSmartContractTransactionMetadataJson returns a JSON-encoded object containing
// the blockchain metadata (e.g. chain_id, block_number, timestamp) of the transaction
// identified by transactionID.
func (client *Client) GetSmartContractTransactionMetadataJson(transactionId string) ([]byte, error) {
	return callUserMan(client, API_GET_SMARTCONTRACT_TRANSACTION_META_JSON, []byte(transactionId))
}

// ListLatestTransactions returns a JSON-encoded object containing a list of the latest
// count transaction IDs sorted by transaction time in descending order (latest first).
func (client *Client) ListLatestTransactions(count int) ([]byte, error) {
	return callUserManV(client, API_LIST_LATEST_TRANSACTION, count)
}
