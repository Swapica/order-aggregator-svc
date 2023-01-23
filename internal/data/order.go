package data

import (
	"database/sql"
	"math/big"
)

type Orders interface {
	Insert(Order) error
	Update(id, srcChain string, state uint8, executedBy *big.Int, matchSwapica *string) error
	Get(id, srcChain string) (*Order, error)
}

type Order struct {
	ID           string `structs:"id" db:"id"`
	SrcChain     string `structs:"src_chain" db:"src_chain"`
	Account      string `structs:"account" db:"account"`
	TokensToSell string `structs:"sell_tokens" db:"sell_tokens"`
	TokensToBuy  string `structs:"buy_tokens" db:"buy_tokens"`
	AmountToSell string `structs:"sell_amount" db:"sell_amount"`
	AmountToBuy  string `structs:"buy_amount" db:"buy_amount"`
	DestChain    string `structs:"dest_chain" db:"dest_chain"`

	State        uint8          `structs:"state" db:"state"`
	ExecutedBy   sql.NullString `structs:"executed_by,omitempty,omitnested" db:"executed_by"`
	MatchSwapica sql.NullString `structs:"match_swapica,omitempty,omitnested" db:"match_swapica"`
}
