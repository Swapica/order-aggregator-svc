package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	matchesTable   = "match_orders m"
	matchesColumns = "m.id,m.src_chain,m.order_id,m.order_chain,m.account,m.sell_token,m.sell_amount,m.state"
)

const (
	orderStateAwaitingFinalization = iota + 2
	orderStateCanceled
	orderStateExecuted
)

type matches struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewMatchOrders(db *pgdb.DB) data.MatchOrders {
	return &matches{
		db:       db,
		selector: squirrel.Select(matchesColumns).From(matchesTable),
		updater:  squirrel.Update(matchesTable),
	}
}

func (q *matches) New() data.MatchOrders {
	return NewMatchOrders(q.db)
}

func (q *matches) Insert(order data.Match) error {
	stmt := squirrel.Insert(matchesTable).SetMap(structs.Map(order))
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to insert match order")
}

func (q *matches) Update(id string, state uint8) error {
	stmt := q.updater.Where(squirrel.Eq{"m.id": id}).Set("m.state", state)
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to update match order")
}

func (q *matches) Get(id string) (*data.Match, error) {
	var res data.Match
	stmt := squirrel.Select("*").From(matchesTable).Where(squirrel.Eq{"m.id": id})
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
	q.selector = page.ApplyTo(q.selector, "m.id")
	return q
}

func (q *matches) FilterByChain(id string) data.MatchOrders {
	return q.filterByCol("m.src_chain", &id)
}

func (q *matches) FilterByAccount(address *string) data.MatchOrders {
	return q.filterByCol("m.account", address)
}

func (q *matches) FilterByState(state *string) data.MatchOrders {
	return q.filterByCol("m.state", state)
}

func (q *matches) FilterExpired(apply *bool) data.MatchOrders {
	if apply == nil || !*apply {
		return q
	}

	q.selector = q.selector.Join(ordersTable + " o ON m.order_id = o.id AND m.order_chain = o.src_chain").Where(
		squirrel.Eq{
			"m.state": orderStateAwaitingFinalization,
			"o.state": []int{orderStateCanceled, orderStateExecuted}}).Where(
		"o.executed_by IS DISTINCT FROM m.id") // works with NULLs better than != or squirrel.NotEq
	return q
}

func (q *matches) filterByCol(column string, value *string) data.MatchOrders {
	if value == nil {
		return q
	}
	q.selector = q.selector.Where(squirrel.Eq{column: value})
	q.updater = q.updater.Where(squirrel.Eq{column: value})
	return q
}
