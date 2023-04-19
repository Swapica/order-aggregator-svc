package responses

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewMatch(m data.Match, srcChain, originChain resources.Key) resources.MatchResponse {
	return resources.MatchResponse{Data: ToMatchResource(m, srcChain, originChain)}
}

func NewMatchList(matches []resources.Match, included []resources.Resource, count int64) resources.MatchListResponse {
	resp := resources.MatchListResponse{Data: matches, Meta: toRawMetaField(count)}
	for i := range included {
		resp.Included.Add(included[i])
	}
	return resp
}

func ToMatchResource(m data.Match, srcChain, originChain resources.Key) resources.Match {
	originKey := resources.NewKeyInt64(m.OriginOrder, resources.ORDER)
	sellKey := resources.NewKeyInt64(m.SellToken, resources.TOKEN)

	return resources.Match{
		Key: resources.NewKeyInt64(m.ID, resources.MATCH_ORDER),
		Attributes: resources.MatchAttributes{
			Creator:       m.Creator,
			AmountToSell:  m.SellAmount,
			MatchId:       m.MatchID,
			OriginOrderId: m.OrderID,
			State:         m.State,
			UseRelayer:    m.UseRelayer,
		},
		Relationships: resources.MatchRelationships{
			OriginChain: resources.Relation{Data: &originChain},
			OriginOrder: resources.Relation{Data: &originKey},
			SrcChain:    resources.Relation{Data: &srcChain},
			TokenToSell: resources.Relation{Data: &sellKey},
		},
	}
}
