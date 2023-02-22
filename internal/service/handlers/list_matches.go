package handlers

import (
	"net/http"

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

	matchesRes := make([]resources.Match, 0, len(matches))
	ordersRes := make([]resources.Order, 0, len(matches))
	orderIDs := make([]int64, 0, len(matches))
	chains := make([]resources.Chain, 0, 2*len(matches))

	for _, m := range matches {
		src := ChainsQ(r).FilterByChainID(m.SrcChain).Get()
		origin := ChainsQ(r).FilterByChainID(m.OrderChain).Get()
		matchesRes = append(matchesRes, responses.ToMatchResource(m, src.Key, origin.Key))

		if req.IncludeSrcChain {
			chains = append(chains, *src)
		}
		if req.IncludeOriginChain {
			chains = append(chains, *origin)
		}
		if req.IncludeOriginOrder {
			orderIDs = append(orderIDs, m.OriginOrder)
		}
	}

	if req.IncludeOriginOrder {
		orders, err := OrdersQ(r).FilterByID(orderIDs...).Select()
		if err != nil {
			Log(r).WithError(err).Error("failed to include orders")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		for _, o := range orders {
			src := ChainsQ(r).FilterByChainID(o.SrcChain).Get()
			dest := ChainsQ(r).FilterByChainID(o.DestChain).Get()
			ordersRes = append(ordersRes, responses.ToOrderResource(o, src.Key, dest.Key))
		}
	}

	resp := responses.NewMatchList(matchesRes, ordersRes, chains, count)
	resp.Links = req.Params.GetLinks(r)
	ape.Render(w, resp)
}
