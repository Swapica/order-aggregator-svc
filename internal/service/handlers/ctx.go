package handlers

import (
	"context"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/internal/data/postgres"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	ordersCtxKey
	matchOrdersCtxKey
	blockCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxOrdersQ(db *pgdb.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ordersCtxKey, postgres.NewOrders(db))
	}
}

func OrdersQ(r *http.Request) data.Orders {
	return r.Context().Value(ordersCtxKey).(data.Orders)
}

func CtxMatchOrdersQ(db *pgdb.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, matchOrdersCtxKey, postgres.NewMatchOrders(db))
	}
}

func MatchOrdersQ(r *http.Request) data.MatchOrders {
	return r.Context().Value(matchOrdersCtxKey).(data.MatchOrders)
}

func CtxBlockQ(db *pgdb.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, blockCtxKey, postgres.NewLastBlock(db))
	}
}

func BlockQ(r *http.Request) data.LastBlock {
	return r.Context().Value(blockCtxKey).(data.LastBlock)
}
