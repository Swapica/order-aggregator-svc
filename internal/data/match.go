package data

import "gitlab.com/distributed_lab/kit/pgdb"

type MatchOrders interface {
	New() MatchOrders
	Insert(Match) error
	Update(id string, state uint8) error
	Get(id string) (*Match, error)
	Select() ([]Match, error)
	Page(*pgdb.CursorPageParams) MatchOrders
	FilterByChain(chainID string) MatchOrders
	FilterByAccount(*string) MatchOrders
	FilterByState(*string) MatchOrders
	FilterExpired(apply bool) MatchOrders
}

type Match struct {
	ID           string `structs:"id" db:"id"`
	SrcChain     string `structs:"src_chain" db:"src_chain"`
	OrderID      string `structs:"order_id" db:"order_id"`
	OrderChain   string `structs:"order_chain" db:"order_chain"`
	Account      string `structs:"account" db:"account"`
	TokenToSell  string `structs:"sell_token" db:"sell_token"`
	AmountToSell string `structs:"sell_amount" db:"sell_amount"`
	State        uint8  `structs:"state" db:"state"`
}
