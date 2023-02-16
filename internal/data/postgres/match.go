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
	matchesTable   = "match_orders"
	matchesColumns = "m.id,m.match_id,m.src_chain,m.order_id,m.order_chain,m.creator,m.sell_token,m.sell_amount,m.state"
)

type matches struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewMatchOrders(db *pgdb.DB) data.MatchOrders {
	return &matches{
		db:       db,
		selector: squirrel.Select(matchesColumns).From(matchesTable + " m"),
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

func (q *matches) Page(page *pgdb.CursorPageParams) data.MatchOrders {
	q.selector = page.ApplyTo(q.selector, "m.id")
	return q
}

func (q *matches) FilterBySupportedChains(chainIDs ...int64) data.MatchOrders {
	condition := squirrel.Eq{"m.src_chain": chainIDs, "m.order_chain": chainIDs}
	q.selector = q.selector.Where(condition)
	q.updater = q.updater.Where(condition)
	return q
}

func (q *matches) FilterByMatchID(id int64) data.MatchOrders {
	return q.filterByCol("match_id", id)
}

func (q *matches) FilterByChain(id *int64) data.MatchOrders {
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

	q.selector = q.selector.Join(ordersTable + " o ON m.order_id = o.order_id AND m.order_chain = o.src_chain").Where(
		squirrel.Eq{
			"m.state": data.StateAwaitingFinalization,
			"o.state": []uint8{data.StateCanceled, data.StateExecuted}}).Where(
		"o.match_id IS DISTINCT FROM m.match_id") // works with NULLs better than != or squirrel.NotEq
	return q
}

func (q *matches) filterByCol(column string, value interface{}) data.MatchOrders {
	if isNilInterface(value) {
		return q
	}

	if _, ok := value.(*string); ok {
		q.selector = q.selector.Where(squirrel.ILike{"m." + column: value})
		q.updater = q.updater.Where(squirrel.ILike{column: value})
		return q
	}

	q.selector = q.selector.Where(squirrel.Eq{"m." + column: value})
	q.updater = q.updater.Where(squirrel.Eq{column: value})
	return q
}
