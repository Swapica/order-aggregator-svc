package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"github.com/Swapica/order-aggregator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ListMatches(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewListMatches(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := MatchOrdersQ(r).
		FilterBySupportedChains(ChainsQ(r).SelectIDs()...).
		FilterBySrcChain(req.FilterSrcChain).
		FilterByCreator(req.FilterCreator).
		FilterByState(req.FilterState).
		FilterExpired(req.FilterExpired)

	matches, err := q.Page(&req.OffsetPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get match orders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	count, err := q.Count()
	if err != nil {
		Log(r).WithError(err).Error("failed to count match orders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var orders []data.Order
	matchesRes := make([]resources.Match, 0, len(matches))
	orderIDs := make([]int64, 0, len(matches))
	tokenIDs := make([]int64, 0, 3*len(matches))
	included := make([]resources.Resource, 0, 4*len(matches))
	includeOrigin := req.IncludeOriginOrder || req.IncludeOriginBuyToken || req.IncludeOriginSellToken

	for _, m := range matches {
		src := ChainsQ(r).FilterByChainID(m.SrcChain).Get()
		origin := ChainsQ(r).FilterByChainID(m.OrderChain).Get()
		matchesRes = append(matchesRes, responses.ToMatchResource(m, src.Key, origin.Key))

		if req.IncludeSrcChain {
			included = append(included, src)
		}
		if req.IncludeOriginChain {
			included = append(included, origin)
		}
		if includeOrigin {
			orderIDs = append(orderIDs, m.OriginOrder)
		}
		if req.IncludeSellToken {
			tokenIDs = append(tokenIDs, m.SellToken)
		}
	}

	if includeOrigin {
		orders, err = OrdersQ(r).FilterByID(orderIDs...).Select()
		if err != nil {
			Log(r).WithError(err).Error("failed to include origin orders")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		for _, o := range orders {
			src := ChainsQ(r).FilterByChainID(o.SrcChain).Get()
			dest := ChainsQ(r).FilterByChainID(o.DestChain).Get()
			res := responses.ToOrderResource(o, src.Key, dest.Key)
			included = append(included, &res)

			if req.IncludeOriginBuyToken {
				tokenIDs = append(tokenIDs, o.BuyToken)
			}
			if req.IncludeOriginSellToken {
				tokenIDs = append(tokenIDs, o.SellToken)
			}
		}
	}

	if req.IncludeSellToken || req.IncludeOriginBuyToken || req.IncludeOriginSellToken {
		tokens, err := TokensQ(r).FilterByID(tokenIDs...).Select()
		if err != nil {
			Log(r).WithError(err).Error("failed to include tokens")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		for _, t := range tokens {
			res := responses.ToTokenResource(t)
			included = append(included, &res)
		}
	}

	resp := responses.NewMatchList(matchesRes, included, count)
	resp.Links = req.Params.GetLinks(r)
	ape.Render(w, resp)
}
