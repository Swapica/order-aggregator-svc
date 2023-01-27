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

type matches struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewMatchOrders(db *pgdb.DB) data.MatchOrders {
	return &matches{
		db:       db,
		selector: squirrel.Select("*").From(matchOrdersTable),
		updater:  squirrel.Update(matchOrdersTable),
	}
}

func (q *matches) New() data.MatchOrders {
	return NewMatchOrders(q.db)
}

func (q *matches) Insert(order data.Match) error {
	stmt := squirrel.Insert(matchOrdersTable).SetMap(structs.Map(order))
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to insert match order")
}

func (q *matches) Update(id string, state uint8) error {
	stmt := q.updater.Where(squirrel.Eq{"id": id}).Set("state", state)
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to update match order")
}

func (q *matches) Get(id string) (*data.Match, error) {
	var res data.Match
	stmt := squirrel.Select("*").From(matchOrdersTable).Where(squirrel.Eq{"id": id})
	err := q.db.Get(&res, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &res, errors.Wrap(err, "failed to get match order")
}

func (q *matches) Select() ([]data.Match, error) {
	var res []data.Match
	err := q.db.Select(&res, q.selector)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return res, errors.Wrap(err, "failed to select match orders")
}

func (q *matches) Page(page *pgdb.CursorPageParams) data.MatchOrders {
	q.selector = page.ApplyTo(q.selector, "id")
	return q
}

func (q *matches) FilterByChain(name string) data.MatchOrders {
	q.selector = q.selector.Where(squirrel.Eq{"src_chain": name})
	q.updater = q.updater.Where(squirrel.Eq{"src_chain": name})
	return q
}
