//
// Copyright 2019 Digital Transaction Limited.
// All Rights Reserved.
//

package parallelcore_client_sdk_go

import (
	pb "github.com/digital-transaction/parallelcore-client-sdk-go/engine_client_proto"
)

/******************************************************************************/
/*                                                                            */
/* Get Blockchain Summary (JSON)                                              */
/*                                                                            */
/******************************************************************************/
func (client *Client) GetBlockchainSummaryJson() ([]byte, error) {
	return callUserMan(client, API_GET_BLOCK_CHAIN_SUMMARY_JSON, make([]byte, 0))
}

/******************************************************************************/
/*                                                                            */
/* Get Blockchain Block Details (JSON)                                        */
/*                                                                            */
/******************************************************************************/
func (client *Client) GetBlockDetailsJson(chainId string, blockId string) ([]byte, error) {
	return callUserManV(client, API_GET_BLOCK_DETAILS_JSON, pb.BlockData{ChainId: chainId, BlockId: blockId})
}

/******************************************************************************/
/*                                                                            */
/* Calculate Blockchain Block Hash                                            */
/*                                                                            */
/******************************************************************************/
func (client *Client) CalculateBlockHash(chainId string, blockId string) ([]byte, error) {
	return callUserManV(client, API_CALCULATE_BLOCK_HASH, pb.BlockData{ChainId: chainId, BlockId: blockId})
}
