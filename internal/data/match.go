package data

import "gitlab.com/distributed_lab/kit/pgdb"

type MatchOrders interface {
	New() MatchOrders
	Insert(Match) (Match, error)
	Update(state uint8) error
	Get() (*Match, error)
	Select() ([]Match, error)
	Count() (int64, error)
	Page(*pgdb.CursorPageParams) MatchOrders
	FilterBySupportedChains(chainIDs ...int64) MatchOrders
	FilterByMatchID(int64) MatchOrders
	FilterBySrcChain(*int64) MatchOrders
	FilterByCreator(*string) MatchOrders
	FilterByState(*uint8) MatchOrders
	FilterExpired(*bool) MatchOrders
}

type Match struct {
	ID         int64  `structs:"-" db:"id"`
	SrcChain   int64  `structs:"src_chain" db:"src_chain"`
	MatchID    int64  `structs:"match_id" db:"match_id"`
	OrderID    int64  `structs:"order_id" db:"order_id"`
	OrderChain int64  `structs:"order_chain" db:"order_chain"`
	Creator    string `structs:"creator" db:"creator"`
	SellToken  string `structs:"sell_token" db:"sell_token"`
	SellAmount string `structs:"sell_amount" db:"sell_amount"`
	State      uint8  `structs:"state" db:"state"`
}
