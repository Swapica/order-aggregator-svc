package responses

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewMatch(m data.Match, srcChain, originChain resources.Key) resources.MatchResponse {
	return resources.MatchResponse{Data: ToMatchResource(m, srcChain, originChain)}
}

func NewMatchList(matches []resources.Match, included []resources.Chain) resources.MatchListResponse {
	resp := resources.MatchListResponse{Data: matches}
	for i := range included {
		resp.Included.Add(&included[i])
	}
	return resp
}

func ToMatchResource(m data.Match, srcChain, originChain resources.Key) resources.Match {
	originOrder := resources.NewKeyInt64(m.OrderID, resources.ORDER)

	return resources.Match{
		Key: resources.NewKeyInt64(m.ID, resources.MATCH_ORDER),
		Attributes: resources.MatchAttributes{
			Creator:      m.Creator,
			AmountToSell: m.SellAmount,
			MatchId:      &m.MatchID,
			State:        m.State,
			TokenToSell:  m.SellToken,
		},
		Relationships: resources.MatchRelationships{
			OriginChain: resources.Relation{Data: &originChain},
			OriginOrder: resources.Relation{Data: &originOrder},
			SrcChain:    resources.Relation{Data: &srcChain},
		},
	}
}
