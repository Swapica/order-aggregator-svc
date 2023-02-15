package responses

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewOrder(o data.Order) resources.OrderResponse {
	return resources.OrderResponse{Data: newOrderResource(o)}
}

func NewOrderList(orders []data.Order, chains []resources.Chain) resources.OrderListResponse {
	var resp resources.OrderListResponse
	resp.Data = make([]resources.Order, len(orders))
	for i, o := range orders {
		resp.Data[i] = newOrderResource(o)
	}
	for i := range chains {
		resp.Included.Add(&chains[i])
	}
	return resp
}

func newOrderResource(o data.Order) resources.Order {
	srcChain := resources.NewKeyInt64(o.SrcChain, resources.CHAIN)
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
			State:        o.State,
			TokenToBuy:   o.BuyToken,
			TokenToSell:  o.SellToken,
		},
		Relationships: resources.OrderRelationships{
			DestinationChain: resources.Relation{Data: &destChain},
			Match:            matchId,
			SrcChain:         resources.Relation{Data: &srcChain},
		},
	}
}
