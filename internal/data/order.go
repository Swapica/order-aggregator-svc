package data

import (
	"database/sql"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Orders interface {
	New() Orders
	Insert(Order) (Order, error)
	Update(state uint8, executedByMatch, matchId *int64, matchSwapica *string) error
	Get() (*Order, error)
	Select() ([]Order, error)
	Count() (int64, error)
	Page(params *pgdb.OffsetPageParams) Orders
	FilterBySupportedChains(chainIDs ...int64) Orders
	FilterByID(ids ...int64) Orders
	FilterByOrderID(ids ...int64) Orders
	FilterBySrcChain(*int64) Orders
	FilterByCreator(*string) Orders
	FilterByTokenToBuy(*string) Orders
	FilterByTokenToSell(*string) Orders
	FilterByDestChain(*int64) Orders
	FilterByState(*uint8) Orders
}

// Order Fields ID and ExecutedByMatch are database-generated properties, any other come from the
// Swapica contract on the SrcChain network.
// MatchID field is preserved for convenience and to reduce computing power and network usage.
type Order struct {
	// ID surrogate key is strongly preferred against PRIMARY KEY (OrderID, SrcChain)
	ID         int64  `structs:"-" db:"id"`
	OrderID    int64  `structs:"order_id" db:"order_id"`
	SrcChain   int64  `structs:"src_chain" db:"src_chain"`
	Creator    string `structs:"creator" db:"creator"`
	SellToken  int64  `structs:"sell_token" db:"sell_token"`
	BuyToken   int64  `structs:"buy_token" db:"buy_token"`
	SellAmount string `structs:"sell_amount" db:"sell_amount"`
	BuyAmount  string `structs:"buy_amount" db:"buy_amount"`
	DestChain  int64  `structs:"dest_chain" db:"dest_chain"`
	State      uint8  `structs:"state" db:"state"`

	// ExecutedByMatch foreign key for match_orders(ID)
	ExecutedByMatch sql.NullInt64  `structs:"executed_by_match,omitempty,omitnested" db:"executed_by_match"`
	MatchID         sql.NullInt64  `structs:"match_id,omitempty,omitnested" db:"match_id"`
	MatchSwapica    sql.NullString `structs:"match_swapica,omitempty,omitnested" db:"match_swapica"`
}
