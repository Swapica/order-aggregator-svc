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

	var matchId *resources.Relation
	if o.MatchID.String != "" {
		matchId = &resources.Relation{
			Data: &resources.Key{
				ID:   o.MatchID.String,
				Type: resources.MATCH_ORDER,
			},
		}
	}

	return resources.Order{
		Key: resources.NewKeyInt64(o.ID, resources.ORDER),
		Attributes: resources.OrderAttributes{
			Creator:      o.Creator,
			AmountToBuy:  o.BuyAmount,
			AmountToSell: o.SellAmount,
			MatchSwapica: matchSwapica,
			OrderId:      &o.OrderID,
			SrcChain:     &o.SrcChain,
			State:        o.State,
			TokenToBuy:   o.BuyToken,
			TokenToSell:  o.SellToken,
		},
		Relationships: resources.OrderRelationships{
			DestinationChain: resources.Relation{Data: &destChain},
			Match:            matchId,
		},
	}
}
