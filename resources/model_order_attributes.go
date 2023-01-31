/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type OrderAttributes struct {
	// Order creator
	Account string `json:"account"`
	// With decimals
	AmountToBuy string `json:"amountToBuy"`
	// With decimals
	AmountToSell string `json:"amountToSell"`
	// Chain ID of the destination network
	DestChain string `json:"destChain"`
	// Match order's ID that allowed to execute the order
	ExecutedBy *string `json:"executedBy,omitempty"`
	// Swapica contract address on the destination network
	MatchSwapica *string `json:"matchSwapica,omitempty"`
	// Order state
	State uint8 `json:"state"`
	// Contract address of the token to buy
	TokenToBuy string `json:"tokenToBuy"`
	// Contract address of the token to sell
	TokenToSell string `json:"tokenToSell"`
}
