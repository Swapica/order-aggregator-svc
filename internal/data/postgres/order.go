package postgres

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/fatih/structs"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const ordersTable = "orders"

type orders struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
	counter  squirrel.SelectBuilder
	updater  squirrel.UpdateBuilder
}

func NewOrders(db *pgdb.DB) data.Orders {
	return &orders{
		db:       db,
		selector: squirrel.Select("*").From(ordersTable),
		counter:  squirrel.Select("count(id)").From(ordersTable),
		updater:  squirrel.Update(ordersTable),
	}
}

func (q *orders) New() data.Orders {
	return NewOrders(q.db)
}

func (q *orders) Insert(order data.Order) (data.Order, error) {
	var res data.Order
	stmt := squirrel.Insert(ordersTable).SetMap(structs.Map(order)).Suffix("RETURNING *")
	err := q.db.Get(&res, stmt)
	return res, errors.Wrap(err, "failed to insert order")
}

func (q *orders) Update(state uint8, matchID *int64, matchSw *string) error {
	stmt := q.updater.SetMap(map[string]interface{}{"state": state, "match_id": matchID, "match_swapica": matchSw})
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to update order")
}

func (q *orders) Get() (*data.Order, error) {
	var res data.Order
	err := q.db.Get(&res, q.selector)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &res, errors.Wrap(err, "failed to get order")
}

func (q *orders) Select() ([]data.Order, error) {
	var res []data.Order
	err := q.db.Select(&res, q.selector)
	return res, errors.Wrap(err, "failed to select orders")
}

func (q *orders) Count() (int64, error) {
	var res struct {
		Count int64 `db:"count"`
	}
	err := q.db.Get(&res, q.counter)
	return res.Count, errors.Wrap(err, "failed to count orders in DB")
}

func (q *orders) Page(page *pgdb.OffsetPageParams) data.Orders {
	// Count() counts all the available records, therefore pagination is not applied to it
	q.selector = page.ApplyTo(q.selector, "id")
	return q
}

func (q *orders) FilterBySupportedChains(chainIDs ...int64) data.Orders {
	return q.filterByCol("src_chain", chainIDs).filterByCol("dest_chain", chainIDs)
}

func (q *orders) FilterByOrderID(ids ...int64) data.Orders {
	return q.filterByCol("order_id", ids)
}

func (q *orders) FilterByCreator(address *string) data.Orders {
	return q.filterByCol("creator", address)
}

func (q *orders) FilterBySrcChain(id *int64) data.Orders {
	return q.filterByCol("src_chain", id)
}

func (q *orders) FilterByTokenToBuy(address *string) data.Orders {
	return q.filterByCol("buy_token", address)
}

func (q *orders) FilterByTokenToSell(address *string) data.Orders {
	return q.filterByCol("sell_token", address)
}

func (q *orders) FilterByDestinationChain(id *int64) data.Orders {
	return q.filterByCol("dest_chain", id)
}

func (q *orders) FilterByState(state *uint8) data.Orders {
	return q.filterByCol("state", state)
}

func (q *orders) filterByCol(column string, value interface{}) *orders {
	if isNilInterface(value) {
		return q
	}

	if _, ok := value.(*string); ok {
		q.selector = q.selector.Where(squirrel.ILike{column: value})
		q.counter = q.counter.Where(squirrel.ILike{column: value})
		q.updater = q.updater.Where(squirrel.ILike{column: value})
		return q
	}

	q.selector = q.selector.Where(squirrel.Eq{column: value})
	q.counter = q.counter.Where(squirrel.Eq{column: value})
	q.updater = q.updater.Where(squirrel.Eq{column: value})
	return q
}
