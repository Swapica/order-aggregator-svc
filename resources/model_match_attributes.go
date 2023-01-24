/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "math/big"

type MatchAttributes struct {
	// **Match order** creator
	Account string `json:"account"`
	// With decimals
	AmountToSell *big.Int `json:"amountToSell"`
	// Chain ID of the order's origin network
	OriginChain *big.Int `json:"originChain"`
	// ID of the order to match
	OriginOrderId *big.Int `json:"originOrderId"`
	// Code name of the **match order's** source chain
	SrcChain string `json:"srcChain"`
	// Order state
	State uint8 `json:"state"`
	// Contract address of the token to sell
	TokensToSell string `json:"tokensToSell"`
}
