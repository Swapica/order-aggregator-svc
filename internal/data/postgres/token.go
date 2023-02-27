package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const tokensTable = "tokens"

type tokens struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func NewTokens(db *pgdb.DB) data.Tokens {
	return &tokens{
		db:       db,
		selector: squirrel.Select("*").From(tokensTable),
	}
}

func (q *tokens) New() data.Tokens {
	return NewTokens(q.db)
}

func (q *tokens) Insert(token data.Token) (data.Token, error) {
	var res data.Token
	stmt := squirrel.Insert(tokensTable).SetMap(structs.Map(token)).Suffix("RETURNING *")
	err := q.db.Get(&res, stmt)
	return res, errors.Wrap(err, "failed to insert token")
}

func (q *tokens) Get() (*data.Token, error) {
	var res data.Token
	err := q.db.Get(&res, q.selector)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &res, errors.Wrap(err, "failed to get token")
}

func (q *tokens) Select() ([]data.Token, error) {
	var res []data.Token
	err := q.db.Select(&res, q.selector)
	return res, errors.Wrap(err, "failed to select tokens")
}

func (q *tokens) FilterByID(ids ...int64) data.Tokens {
	return q.filterByCol("id", ids)
}

func (q *tokens) FilterByAddress(a string) data.Tokens {
	return q.filterByCol("address", a)
}

func (q *tokens) FilterByChain(chainID int64) data.Tokens {
	return q.filterByCol("src_chain", chainID)
}

func (q *tokens) filterByCol(column string, value interface{}) *tokens {
	if isNilInterface(value) {
		return q
	}

	if _, ok := value.(*string); ok {
		q.selector = q.selector.Where(squirrel.ILike{column: value})
		return q
	}

	q.selector = q.selector.Where(squirrel.Eq{column: value})
	return q
}
