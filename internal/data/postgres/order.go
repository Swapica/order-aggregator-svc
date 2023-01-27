package postgres

import (
	"database/sql"
	"math/big"

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
	updater  squirrel.UpdateBuilder
}

func NewOrders(db *pgdb.DB) data.Orders {
	return &orders{
		db:       db,
		selector: squirrel.Select("*").From(ordersTable),
		updater:  squirrel.Update(ordersTable),
	}
}

func (q *orders) New() data.Orders {
	return NewOrders(q.db)
}

func (q *orders) Insert(order data.Order) error {
	stmt := squirrel.Insert(ordersTable).SetMap(structs.Map(order))
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to insert order")
}

func (q *orders) Update(id string, state uint8, execBy *big.Int, matchSw *string) error {
	updMap := map[string]interface{}{"state": state, "executed_by": nil, "match_swapica": matchSw}
	if execBy != nil {
		updMap["executed_by"] = execBy.String()
	}

	stmt := q.updater.Where(squirrel.Eq{"id": id}).SetMap(updMap)
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to update order")
}

func (q *orders) Get(id string) (*data.Order, error) {
	var res data.Order
	err := q.db.Get(&res, q.selector.Where(squirrel.Eq{"id": id}))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &res, errors.Wrap(err, "failed to get order")
}

func (q *orders) Select() ([]data.Order, error) {
	var res []data.Order
	err := q.db.Select(&res, q.selector)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return res, errors.Wrap(err, "failed to select orders")
}

func (q *orders) Page(page *pgdb.CursorPageParams) data.Orders {
	q.selector = page.ApplyTo(q.selector, "id")
	return q
}

func (q *orders) FilterByChain(name string) data.Orders {
	return q.filterByCol("src_chain", &name)
}

func (q *orders) FilterByTokenToBuy(address *string) data.Orders {
	return q.filterByCol("buy_token", address)
}

func (q *orders) FilterByTokenToSell(address *string) data.Orders {
	return q.filterByCol("sell_token", address)
}

func (q *orders) FilterByAccount(address *string) data.Orders {
	return q.filterByCol("account", address)
}

func (q *orders) FilterByState(state *string) data.Orders {
	return q.filterByCol("state", state)
}

func (q *orders) filterByCol(column string, value *string) data.Orders {
	if value == nil {
		return q
	}
	q.selector = q.selector.Where(squirrel.Eq{column: value})
	q.updater = q.updater.Where(squirrel.Eq{column: value})
	return q
}
