package responses

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewMatchListResponse(matches []data.Match) resources.MatchListResponse {
	list := make([]resources.Match, len(matches))
	for i, o := range matches {
		list[i] = newMatchResource(o)
	}
	return resources.MatchListResponse{Data: list}
}

func newMatchResource(o data.Match) resources.Match {
	return resources.Match{
		Key: resources.Key{
			ID:   o.ID,
			Type: resources.MATCH_ORDER,
		},
		Attributes: resources.MatchAttributes{
			Account:      o.Account,
			AmountToSell: parseBig(o.AmountToSell, "amountToSell"),
			OriginChain:  parseBig(o.OrderChain, "originChain"),
			State:        o.State,
			TokenToSell:  o.TokenToSell,
		},
		Relationships: resources.MatchRelationships{
			OriginOrder: resources.Relation{
				Data: &resources.Key{
					ID:   o.OrderID,
					Type: resources.ORDER,
				},
			},
		},
	}
}
