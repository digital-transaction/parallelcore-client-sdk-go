//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

type InfoListData struct {
	AllDomains bool   `json:"allDomains"`
	DomainName string `json:"domainName"`
}

type UserManData struct {
	Action string `json:"action"`
	Data   []byte `json:"data"`
}

type SysManData struct {
	Action string `json:"action"`
	Data   []byte `json:"data"`
}

type UserData struct {
	ID         string `json:"clientId"`
	Credential string `json:"clientCredential"`
	Roles      string `json:"clientRoles"`
	DomainName string `json:"clientDomainName"`
}

type UserDomainData struct {
	ID         string `json:"clientId"`
	DomainName string `json:"clientDomainName"`
}

type UserAccessData struct {
	ID                string `json:"clientId"`
	SmartContractName string `json:"scName"`
	DomainName        string `json:"domainName"`
}

type UserFullData struct {
	ID             string   `json:"clientId"`
	Roles          []string `json:"roles"`
	AccessList     []string `json:"accessList"`
	Domains        []string `json:"domains"`
	ManagedDomains []string `json:"managedDomains"`
}

type UserFullDataWrapper struct {
	ID   string       `json:"clientId"`
	Data UserFullData `json:"data"`
}

type ApiAccessControlData struct {
	Operation string `json:"operation"`
	ApiName   string `json:"api"`
	Options   []byte `json:"options"`
}

type BlockData struct {
	ChainId string `json:"chainId"`
	BlockId string `json:"blockId"`
}

type ForgetGroup struct {
	TxIds []string `json:"tx_ids"`
}

type RequestForgetParams struct {
	TxIds []string `json:"tx_ids"`
}

type ApproveForgetParams struct {
	RequestTxId string `json:"request_tx_id"`
}

type CommitForgetParams struct {
	RequestTxId   string   `json:"request_tx_id"`
	ApprovalTxIds []string `json:"approval_tx_ids"`
}

type ForgetReport struct {
	Deleted        []string `json:"deleted"`
	AlreadyDeleted []string `json:"already_deleted"`
	NotFound       []string `json:"not_found"`
	CommitTxId     string   `json:"commit_tx_id"`
}

type SmartContractData struct {
	Name        string `json:"scName"`
	FileContent []byte `json:"file-content"`
	DomainName  string `json:"domainName"`
	InitArgs    string `json:"init-args"`
}

//types of api access control payloads
type GetSmartContractTransactionOptions struct {
	ClientId      string `json:"client-id"`
	SmartContract string `json:"smart-contract"`
}
