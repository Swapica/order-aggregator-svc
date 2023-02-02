package data

import (
	"database/sql"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Orders interface {
	New() Orders
	Insert(Order) error
	Update(state uint8, executedBy *int64, matchSwapica *string) error
	Get() (*Order, error)
	Select() ([]Order, error)
	Page(*pgdb.CursorPageParams) Orders
	FilterByOrderID(int64) Orders
	FilterByChain(*int64) Orders
	FilterByTokenToBuy(*string) Orders
	FilterByTokenToSell(*string) Orders
	FilterByAccount(*string) Orders
	FilterByState(*uint8) Orders
}

type Order struct {
	ID           int64  `structs:"-" db:"id"`
	SrcChain     int64  `structs:"src_chain" db:"src_chain"`
	OrderID      int64  `structs:"order_id" db:"order_id"`
	Account      string `structs:"account" db:"account"`
	TokenToSell  string `structs:"sell_token" db:"sell_token"`
	TokenToBuy   string `structs:"buy_token" db:"buy_token"`
	AmountToSell string `structs:"sell_amount" db:"sell_amount"`
	AmountToBuy  string `structs:"buy_amount" db:"buy_amount"`
	DestChain    int64  `structs:"dest_chain" db:"dest_chain"`

	State        uint8          `structs:"state" db:"state"`
	ExecutedBy   sql.NullString `structs:"executed_by,omitempty,omitnested" db:"executed_by"`
	MatchSwapica sql.NullString `structs:"match_swapica,omitempty,omitnested" db:"match_swapica"`
}
