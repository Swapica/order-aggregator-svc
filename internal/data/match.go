package data

import "gitlab.com/distributed_lab/kit/pgdb"

type MatchOrders interface {
	New() MatchOrders
	Insert(Match) error
	Update(id string, state uint8) error
	Get(id string) (*Match, error)
	Select() ([]Match, error)
	Page(*pgdb.CursorPageParams) MatchOrders
	FilterByChain(name string) MatchOrders
}

type Match struct {
	ID            string `structs:"id" db:"id"`
	SrcChain      string `structs:"src_chain" db:"src_chain"`
	OriginOrderId string `structs:"origin_order_id" db:"origin_order_id"`
	Account       string `structs:"account" db:"account"`
	TokenToSell   string `structs:"sell_token" db:"sell_token"`
	AmountToSell  string `structs:"sell_amount" db:"sell_amount"`
	OriginChain   string `structs:"origin_chain" db:"origin_chain"`
	State         uint8  `structs:"state" db:"state"`
}
