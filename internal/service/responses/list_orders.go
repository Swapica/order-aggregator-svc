package responses

import (
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
	destChain := resources.NewKeyInt64(o.DestChain, resources.CHAIN)
	var matchSwapica *string
	if m := o.MatchSwapica.String; m != "" {
		matchSwapica = &m
	}

	var executedBy *resources.Relation
	if o.ExecutedBy.String != "" {
		executedBy = &resources.Relation{
			Data: &resources.Key{
				ID:   o.ExecutedBy.String,
				Type: resources.MATCH_ORDER,
			},
		}
	}

	return resources.Order{
		Key: resources.NewKeyInt64(o.ID, resources.ORDER),
		Attributes: resources.OrderAttributes{
			Account:      o.Account,
			AmountToBuy:  o.AmountToBuy,
			AmountToSell: o.AmountToSell,
			MatchSwapica: matchSwapica,
			State:        o.State,
			TokenToBuy:   o.TokenToBuy,
			TokenToSell:  o.TokenToSell,
		},
		Relationships: resources.OrderRelationships{
			DestChain:  resources.Relation{Data: &destChain},
			ExecutedBy: executedBy,
		},
	}
}
