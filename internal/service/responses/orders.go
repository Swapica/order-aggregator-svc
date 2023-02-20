package responses

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewOrder(o data.Order, srcChain, destChain resources.Key) resources.OrderResponse {
	return resources.OrderResponse{Data: ToOrderResource(o, srcChain, destChain)}
}

func NewOrderList(orders []resources.Order, included []resources.Chain, count int64) resources.OrderListResponse {
	resp := resources.OrderListResponse{Data: orders, Meta: toRawMetaField(count)}
	for i := range included {
		resp.Included.Add(&included[i])
	}
	return resp
}

func ToOrderResource(o data.Order, srcChain, destChain resources.Key) resources.Order {
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
