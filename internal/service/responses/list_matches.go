package responses

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

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

func newMatchResource(o data.Match) resources.Match {
	srcChain := resources.NewKeyInt64(o.SrcChain, resources.CHAIN)
	originChain := resources.NewKeyInt64(o.OrderChain, resources.CHAIN)
	originOrder := resources.NewKeyInt64(o.OrderID, resources.ORDER)

	return resources.Match{
		Key: resources.NewKeyInt64(o.ID, resources.MATCH_ORDER),
		Attributes: resources.MatchAttributes{
			Creator:      o.Creator,
			AmountToSell: o.SellAmount,
			MatchId:      &o.MatchID,
			State:        o.State,
			TokenToSell:  o.SellToken,
		},
		Relationships: resources.MatchRelationships{
			OriginChain: resources.Relation{Data: &originChain},
			OriginOrder: resources.Relation{Data: &originOrder},
			SrcChain:    resources.Relation{Data: &srcChain},
		},
	}
}
