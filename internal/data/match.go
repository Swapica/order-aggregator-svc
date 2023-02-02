package data

import "gitlab.com/distributed_lab/kit/pgdb"

type MatchOrders interface {
	New() MatchOrders
	Insert(Match) error
	Update(state uint8) error
	Get() (*Match, error)
	Select() ([]Match, error)
	Page(*pgdb.CursorPageParams) MatchOrders
	FilterByMatchID(int64) MatchOrders
	FilterByChain(*int64) MatchOrders
	FilterByAccount(*string) MatchOrders
	FilterByState(*uint8) MatchOrders
	FilterExpired(*bool) MatchOrders
}

type Match struct {
	ID           int64  `structs:"-" db:"id"`
	SrcChain     int64  `structs:"src_chain" db:"src_chain"`
	MatchID      int64  `structs:"match_id" db:"match_id"`
	OrderID      int64  `structs:"order_id" db:"order_id"`
	OrderChain   int64  `structs:"order_chain" db:"order_chain"`
	Account      string `structs:"account" db:"account"`
	TokenToSell  string `structs:"sell_token" db:"sell_token"`
	AmountToSell string `structs:"sell_amount" db:"sell_amount"`
	State        uint8  `structs:"state" db:"state"`
}
