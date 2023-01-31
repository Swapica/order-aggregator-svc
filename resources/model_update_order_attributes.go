/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UpdateOrderAttributes struct {
	// Match order's ID that allowed to execute the order
	ExecutedBy *string `json:"executedBy,omitempty"`
	// Swapica contract address on the destination network
	MatchSwapica *string `json:"matchSwapica,omitempty"`
	// New order state
	State uint8 `json:"state"`
}
