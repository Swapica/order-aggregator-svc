/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "math/big"

type UpdateOrderAttributes struct {
	// Match order's ID that allowed to execute the order
	ExecutedBy *big.Int `json:"executedBy,omitempty"`
	// Swapica contract address on the destination network
	MatchSwapica *string `json:"matchSwapica,omitempty"`
	// Code name of the order's source chain. **Not to be updated**
	SrcChain string `json:"srcChain"`
	// New order state
	State uint8 `json:"state"`
}
