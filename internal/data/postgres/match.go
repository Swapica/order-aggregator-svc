package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const matchOrdersTable = "match_orders"

type match struct {
	db *pgdb.DB
}

func NewMatchOrders(db *pgdb.DB) data.MatchOrders {
	return match{db: db}
}

func (q match) Insert(order data.Match) error {
	stmt := squirrel.Insert(matchOrdersTable).SetMap(structs.Map(order))
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to insert match order")
}

func (q match) Update(id, chain string, state uint8) error {
	stmt := squirrel.Update(matchOrdersTable).Set("state", state).
		Where(squirrel.Eq{"id": id, "src_chain": chain})
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to update match order")
}

func (q match) Get(id, chain string) (*data.Match, error) {
	var res data.Match
	stmt := squirrel.Select("*").From(matchOrdersTable).Where(squirrel.Eq{"id": id, "src_chain": chain})
	err := q.db.Get(&res, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &res, errors.Wrap(err, "failed to get match order")
}
