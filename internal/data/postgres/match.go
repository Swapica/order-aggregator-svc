package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const matchesTable = "match_orders"

type matches struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	counter  squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewMatchOrders(db *pgdb.DB) data.MatchOrders {
	return &matches{
		db:       db,
		selector: squirrel.Select("m.*").From(matchesTable + " m"),
		counter:  squirrel.Select("count(m.id)").From(matchesTable + " m"),
		updater:  squirrel.Update(matchesTable),
	}
}

func (q *matches) New() data.MatchOrders {
	return NewMatchOrders(q.db)
}

func (q *matches) Insert(order data.Match) (data.Match, error) {
	var res data.Match
	stmt := squirrel.Insert(matchesTable).SetMap(structs.Map(order)).Suffix("RETURNING *")
	err := q.db.Get(&res, stmt)
	return res, errors.Wrap(err, "failed to insert match order")
}

func (q *matches) Update(state uint8) error {
	// update is not supported in FilterExpired, therefore no table alias is needed
	stmt := q.updater.Set("state", state)
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to update match order")
}

func (q *matches) Get() (*data.Match, error) {
	var res data.Match
	err := q.db.Get(&res, q.selector)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &res, errors.Wrap(err, "failed to get match order")
}

func (q *matches) Select() ([]data.Match, error) {
	var res []data.Match
	err := q.db.Select(&res, q.selector)
	return res, errors.Wrap(err, "failed to select match orders")
}

func (q *matches) Count() (int64, error) {
	var res struct {
		Count int64 `db:"count"`
	}
	err := q.db.Get(&res, q.counter)
	return res.Count, errors.Wrap(err, "failed to count match orders in DB")
}

func (q *matches) Page(page *pgdb.OffsetPageParams) data.MatchOrders {
	q.selector = page.ApplyTo(q.selector, "m.id")
	return q
}

func (q *matches) FilterBySupportedChains(chainIDs ...int64) data.MatchOrders {
	condition := squirrel.Eq{"m.src_chain": chainIDs, "m.order_chain": chainIDs}
	q.selector = q.selector.Where(condition)
	q.counter = q.counter.Where(condition)
	q.updater = q.updater.Where(condition)
	return q
}

func (q *matches) FilterByMatchID(id int64) data.MatchOrders {
	return q.filterByCol("match_id", id)
}

func (q *matches) FilterBySrcChain(id *int64) data.MatchOrders {
	return q.filterByCol("src_chain", id)
}

func (q *matches) FilterByCreator(address *string) data.MatchOrders {
	return q.filterByCol("creator", address)
}

func (q *matches) FilterByState(state *uint8) data.MatchOrders {
	return q.filterByCol("state", state)
}

func (q *matches) FilterExpired(apply *bool) data.MatchOrders {
	if apply == nil || !*apply {
		return q
	}

	join := ordersTable + " o ON m.order_id = o.order_id AND m.order_chain = o.src_chain"
	states := squirrel.Eq{
		"m.state": data.StateAwaitingFinalization,
		"o.state": [2]uint8{data.StateCanceled, data.StateExecuted},
	}
	distinct := "o.match_id IS DISTINCT FROM m.match_id" // works with NULLs better than != or squirrel.NotEq

	q.selector = q.selector.Join(join).Where(states).Where(distinct)
	q.counter = q.counter.Join(join).Where(states).Where(distinct)
	return q
}

func (q *matches) filterByCol(column string, value interface{}) data.MatchOrders {
	if isNilInterface(value) {
		return q
	}

	if _, ok := value.(*string); ok {
		q.selector = q.selector.Where(squirrel.ILike{"m." + column: value})
		q.counter = q.counter.Where(squirrel.ILike{"m." + column: value})
		q.updater = q.updater.Where(squirrel.ILike{column: value})
		return q
	}

	q.selector = q.selector.Where(squirrel.Eq{"m." + column: value})
	q.counter = q.counter.Where(squirrel.Eq{"m." + column: value})
	q.updater = q.updater.Where(squirrel.Eq{column: value})
	return q
}
