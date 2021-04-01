//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

// GetBlockchainSummaryJson returns a JSON-encoded object with information about
// all blockchain(s) in the ParallelChain network.
//
// The returned object has fields:
//  - chains Array of Object:
//    {
//    + sealed_block_count int
//    + chain_id string
//    + machine_id string
//    + network_address string
//    + pcore_id string
//    + tags Object:
//      > name string
//    + last_block Object:
//      > block_number int
//      > chunkset_count int
//      > creation_timestamp int
//      > hash string
//      > prev_hash string
//      > status
//    }
func (client *Client) GetBlockchainSummaryJson() ([]byte, error) {
	return callUserMan(client, API_GET_BLOCK_CHAIN_SUMMARY_JSON, make([]byte, 0))
}

// GetBlockDetailsJson returns a JSON-encoded object with information about the block
// identified by chainID and blockID.
//
// The returned object has fields:
//  - block string
//  - block_number int
//  - chain_id string
//  - chunkset_count int
//  - creation_timestamp int
//  - hash string
//  - machine_id string
//  - network_address string
//  - pcore_id string
//  - prev_hash string
//  - status int
func (client *Client) GetBlockDetailsJson(chainID string, blockID string) ([]byte, error) {
	return callUserManV(client, API_GET_BLOCK_DETAILS_JSON, BlockData{ChainId: chainID, BlockId: blockID})
}

// CalculateBlockHash returns a string that is the of the block identified by chainID
// and blockID.
func (client *Client) CalculateBlockHash(chainID string, blockID string) ([]byte, error) {
	return callUserManV(client, API_CALCULATE_BLOCK_HASH, BlockData{ChainId: chainID, BlockId: blockID})
}
