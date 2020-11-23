//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

const (
	DOMAIN_DEFAULT = "default"

	E_FUNC_X_OUTPUT_DECODE_ERROR_X     = "CLIENT: %s: Output decoding Error (%w)"
	E_FUNC_X_ERROR_X                   = "CLIENT: %s: %w"
	FMT_FUNC_X_INPUT_ENCODE_ERROR_X    = "CLIENT: %s: Input encoding Error (%w)"
	FMT_FUNC_X_TASK_ENCODE_ERROR_X     = "CLIENT: %s: Task encoding Error (%w)"
	FMT_FUNC_X_RESPONSE_ERROR_X        = "CLIENT: %s: Response Error (%w)"
	FMT_FUNC_X_OPEN_CONNECTION_ERROR_X = "CLIENT: %s: Failed to open a connection Error (%w)"
	FMT_FUNC_X_NO_SUCCESS_MSG          = "CLIENT: %s: Did not receive expected success message"

	// User Functions
	API_UPDATE_SELF_CREDENTIAL                  = "updateSelfCredential"
	API_LIST_INVOKABLE_SC                       = "ListInvokableSC"
	API_LIST_LATEST_TRANSACTION                 = "ListLatestTransactions"
	API_GET_BLOCK_CHAIN_SUMMARY_JSON            = "GetBlockchainSummaryJson"
	API_GET_BLOCK_DETAILS_JSON                  = "GetBlockDetailsJson"
	API_GET_SMARTCONTRACT_TRANSACTION_JSON      = "GetSmartContractTransactionJson"
	API_GET_SMARTCONTRACT_TRANSACTION_META_JSON = "GetSmartContractTransactionMetadataJson"
	API_CALCULATE_BLOCK_HASH                    = "CalculateBlockHash"

	// System Calls
	API_CREATE_DOMAIN             = "createDomain"
	API_LIST_DOMAIN               = "listDomain"
	API_LIST_MANAGED_DOMAINS      = "listManagedDomains"
	API_GRANT_DOMAIN_ADMIN        = "grantDomainAdmin"
	API_REVOKE_DOMAIN_ADMIN       = "revokeDomainAdmin"
	API_CREATE_CLIENT             = "createClient"
	API_UPDATE_CLIENT             = "updateClient"
	API_REMOVE_CLIENT             = "removeClient"
	API_LIST_CLIENT               = "listClient"
	API_LIST_CLIENTS              = "listClients"
	API_GRANT_ACCESS              = "grantAccess"
	API_REVOKE_ACCESS             = "revokeAccess"
	API_ACCESS                    = "apiAccess"
	API_CHECK_API_ACCESS          = "CheckApiAccess"
	API_MANAGE_API_ACCESS         = "ManageApiAccess"
	API_REGISTER_SMARTCONTRACT    = "registerSmartContract"
	API_LIST_SMARTCONTRACT        = "listSmartContract"
	API_LIST_SMARTCONTRACTS       = "listSmartContracts"
	API_LIST_DOMAIN_SMARTCONTRACT = "listDomainSmartContract"
	API_REGISTER_EVENT_LISTENER   = "registerEventListener"
	API_SYS_MAN                   = "SysMan"
	API_USER_MAN                  = "UserMan"

	// Right To Forget
	API_REQUEST_FORGET     = "requestForget"
	API_APPROVE_FORGET     = "approveForget"
	API_COMMIT_FORGET      = "commitForget"
	API_LIST_FORGET_GROUPS = "listForgetGroups"
)
