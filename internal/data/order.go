package data

import (
	"database/sql"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Orders interface {
	New() Orders
	Insert(Order) error
	Update(state uint8, matchID *int64, matchSwapica *string) error
	Get() (*Order, error)
	Select() ([]Order, error)
	Page(*pgdb.CursorPageParams) Orders
	FilterByOrderID(int64) Orders
	FilterByChain(*int64) Orders
	FilterByCreator(*string) Orders
	FilterByTokenToBuy(*string) Orders
	FilterByTokenToSell(*string) Orders
	FilterByDestinationChain(*int64) Orders
	FilterByState(*uint8) Orders
}

type Order struct {
	ID         int64  `structs:"-" db:"id"`
	SrcChain   int64  `structs:"src_chain" db:"src_chain"`
	OrderID    int64  `structs:"order_id" db:"order_id"`
	Creator    string `structs:"creator" db:"creator"`
	SellToken  string `structs:"sell_token" db:"sell_token"`
	BuyToken   string `structs:"buy_token" db:"buy_token"`
	SellAmount string `structs:"sell_amount" db:"sell_amount"`
	BuyAmount  string `structs:"buy_amount" db:"buy_amount"`
	DestChain  int64  `structs:"dest_chain" db:"dest_chain"`

	State        uint8          `structs:"state" db:"state"`
	MatchID      sql.NullString `structs:"match_id,omitempty,omitnested" db:"match_id"`
	MatchSwapica sql.NullString `structs:"match_swapica,omitempty,omitnested" db:"match_swapica"`
}
