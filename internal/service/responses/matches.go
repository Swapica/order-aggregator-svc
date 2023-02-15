package responses

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewMatch(m data.Match) resources.MatchResponse {
	return resources.MatchResponse{Data: newMatchResource(m)}
}

func NewMatchList(matches []data.Match, chains []resources.Chain) resources.MatchListResponse {
	var resp resources.MatchListResponse
	resp.Data = make([]resources.Match, len(matches))
	for i, o := range matches {
		resp.Data[i] = newMatchResource(o)
	}
	for i := range chains {
		resp.Included.Add(&chains[i])
	}
	return resp
}

func newMatchResource(m data.Match) resources.Match {
	srcChain := resources.NewKeyInt64(m.SrcChain, resources.CHAIN)
	originChain := resources.NewKeyInt64(m.OrderChain, resources.CHAIN)
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
