package data

type Tokens interface {
	New() Tokens
	Insert(Token) (Token, error)
	Get() (*Token, error)
	Select() ([]Token, error)
	FilterByID(ids ...int64) Tokens
	FilterByAddress(string) Tokens
	FilterByChain(int64) Tokens
}

// Token ID field is database-generated property, any other come from the Address contract on the SrcChain network
type Token struct {
	ID       int64  `structs:"-" db:"id"`
	Address  string `structs:"address" db:"address"`
	SrcChain int64  `structs:"src_chain" db:"src_chain"`
	Name     string `structs:"name" db:"name"`
	Symbol   string `structs:"symbol" db:"symbol"`
	Decimals uint8  `structs:"decimals" db:"decimals"`
}
