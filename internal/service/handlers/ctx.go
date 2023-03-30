package handlers

import (
	"context"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/internal/ws"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	ordersCtxKey
	matchOrdersCtxKey
	blockCtxKey
	chainsCtxKey
	tokensCtxKey
	webSocketCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxOrdersQ(q data.Orders) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ordersCtxKey, q)
	}
}

func OrdersQ(r *http.Request) data.Orders {
	return r.Context().Value(ordersCtxKey).(data.Orders).New()
}

func CtxMatchOrdersQ(q data.MatchOrders) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, matchOrdersCtxKey, q)
	}
}

func MatchOrdersQ(r *http.Request) data.MatchOrders {
	return r.Context().Value(matchOrdersCtxKey).(data.MatchOrders).New()
}

func CtxBlockQ(q data.LastBlock) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, blockCtxKey, q)
	}
}

func BlockQ(r *http.Request) data.LastBlock {
	return r.Context().Value(blockCtxKey).(data.LastBlock)
}

func CtxChainsQ(q data.Chains) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, chainsCtxKey, q)
	}
}

func ChainsQ(r *http.Request) data.Chains {
	return r.Context().Value(chainsCtxKey).(data.Chains).New()
}

func CtxTokensQ(q data.Tokens) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, tokensCtxKey, q)
	}
}

func TokensQ(r *http.Request) data.Tokens {
	return r.Context().Value(tokensCtxKey).(data.Tokens).New()
}

func WebSocket(r *http.Request) *ws.Hub {
	return r.Context().Value(webSocketCtxKey).(*ws.Hub)
}

func CtxWebSocket(entry *ws.Hub) func(ctx context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, webSocketCtxKey, entry)
	}
}
