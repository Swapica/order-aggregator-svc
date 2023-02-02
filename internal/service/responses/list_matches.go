package responses

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewMatchList(matches []data.Match) resources.MatchListResponse {
	list := make([]resources.Match, len(matches))
	for i, o := range matches {
		list[i] = newMatchResource(o)
	}
	return resources.MatchListResponse{Data: list}
}

func newMatchResource(o data.Match) resources.Match {
	originOrder := resources.NewKeyInt64(o.OrderID, resources.ORDER)
	originChain := resources.NewKeyInt64(o.OrderChain, resources.CHAIN)

	return resources.Match{
		Key: resources.NewKeyInt64(o.ID, resources.MATCH_ORDER),
		Attributes: resources.MatchAttributes{
			Account:      o.Account,
			AmountToSell: o.AmountToSell,
			MatchId:      &o.MatchID,
			SrcChain:     &o.SrcChain,
			State:        o.State,
			TokenToSell:  o.TokenToSell,
		},
		Relationships: resources.MatchRelationships{
			OriginChain: resources.Relation{
				Data: &originChain,
			},
			OriginOrder: resources.Relation{
				Data: &originOrder,
			},
		},
	}
}
