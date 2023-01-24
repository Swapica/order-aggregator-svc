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
	db *pgdb.DB
}

func NewOrders(db *pgdb.DB) data.Orders {
	return orders{db: db}
}

func (q orders) Insert(order data.Order) error {
	stmt := squirrel.Insert(ordersTable).SetMap(structs.Map(order))
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to insert order")
}

func (q orders) Update(id, chain string, state uint8, execBy *big.Int, matchSw *string) error {
	updMap := map[string]interface{}{"state": state, "executed_by": nil, "match_swapica": matchSw}
	if execBy != nil {
		updMap["executed_by"] = execBy.String()
	}

	stmt := squirrel.Update(ordersTable).SetMap(updMap).Where(squirrel.Eq{"id": id, "src_chain": chain})
	err := q.db.Exec(stmt)
	return errors.Wrap(err, "failed to update order")
}

func (q orders) Get(id, chain string) (*data.Order, error) {
	var res data.Order
	stmt := squirrel.Select("*").From(ordersTable).Where(squirrel.Eq{"id": id, "src_chain": chain})
	err := q.db.Get(&res, stmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &res, errors.Wrap(err, "failed to get order")
}
