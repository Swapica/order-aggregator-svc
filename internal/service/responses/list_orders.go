package responses

import (
	"database/sql"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewOrderList(orders []data.Order) resources.OrderListResponse {
	list := make([]resources.Order, len(orders))
	for i, o := range orders {
		list[i] = newOrderResource(o)
	}
	return resources.OrderListResponse{Data: list}
}

func newOrderResource(o data.Order) resources.Order {
	return resources.Order{
		Key: resources.Key{
			ID:   o.ID,
			Type: resources.ORDER,
		},
		Attributes: resources.OrderAttributes{
			Account:      o.Account,
			AmountToBuy:  o.AmountToBuy,
			AmountToSell: o.AmountToSell,
			DestChain:    o.DestChain,
			ExecutedBy:   nullStringToPtr(o.ExecutedBy),
			MatchSwapica: nullStringToPtr(o.MatchSwapica),
			State:        o.State,
			TokenToBuy:   o.TokenToBuy,
			TokenToSell:  o.TokenToSell,
		},
	}
}

func nullStringToPtr(s sql.NullString) *string {
	if s.String != "" {
		return &s.String
	}
	return nil
}
