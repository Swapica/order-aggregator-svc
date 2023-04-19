package data

import "gitlab.com/distributed_lab/kit/pgdb"

type MatchOrders interface {
	New() MatchOrders
	Insert(Match) (Match, error)
	Update(state uint8) error
	Get() (*Match, error)
	Select() ([]Match, error)
	Count() (int64, error)
	Page(params *pgdb.OffsetPageParams) MatchOrders
	FilterBySupportedChains(chainIDs ...int64) MatchOrders
	FilterByMatchID(int64) MatchOrders
	FilterBySrcChain(*int64) MatchOrders
	FilterByCreator(*string) MatchOrders
	FilterByState(*uint8) MatchOrders
	FilterExpired(*bool) MatchOrders
	FilterClaimable(creator string, srcChain *int64) MatchOrders
	FilterByUseRelayer(*bool) MatchOrders
}

// Match Fields ID and OriginOrder are database-generated properties, any other come from the
// Swapica contract on the SrcChain network
// SrcChain, OrderID, OrderChain, SellToken, SellAmount fields are preserved for convenience
// and to reduce computing power and network usage.
type Match struct {
	// ID surrogate key is strongly preferred against PRIMARY KEY (MatchID, SrcChain)
	ID       int64 `structs:"-" db:"id"`
	MatchID  int64 `structs:"match_id" db:"match_id"`
	SrcChain int64 `structs:"src_chain" db:"src_chain"`
	// OriginOrder foreign key for orders(ID)
	OriginOrder int64  `structs:"origin_order" db:"origin_order"`
	OrderID     int64  `structs:"order_id" db:"order_id"`
	OrderChain  int64  `structs:"order_chain" db:"order_chain"`
	Creator     string `structs:"creator" db:"creator"`
	SellToken   int64  `structs:"sell_token" db:"sell_token"`
	SellAmount  string `structs:"sell_amount" db:"sell_amount"`
	State       uint8  `structs:"state" db:"state"`
	UseRelayer  bool   `structs:"use_relayer" db:"use_relayer"`
}
