//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

/******************************************************************************/
/*                                                                            */
/* Get Smart Contract Transactiono (JSON)                                     */
/*                                                                            */
/******************************************************************************/
func (client *Client) GetSmartContractTransactionJson(transactionId string) ([]byte, error) {
	return callUserMan(client, API_GET_SMARTCONTRACT_TRANSACTION_JSON, []byte(transactionId))
}

/******************************************************************************/
/*                                                                            */
/* Get Smart Contract Transactiono (Meta-Data JSON)                           */
/*                                                                            */
/******************************************************************************/
func (client *Client) GetSmartContractTransactionMetadataJson(transactionId string) ([]byte, error) {
	return callUserMan(client, API_GET_SMARTCONTRACT_TRANSACTION_META_JSON, []byte(transactionId))
}

/******************************************************************************/
/*                                                                            */
/* List Latest Transaction                                                    */
/*                                                                            */
/******************************************************************************/
func (client *Client) ListLatestTransactions(count int) ([]byte, error) {
	return callUserManV(client, API_LIST_LATEST_TRANSACTION, count)
}
