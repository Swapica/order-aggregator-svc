package handlers

import (
	"net/http"
	"strconv"

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

	matches, err := MatchOrdersQ(r).
		FilterBySupportedChains(ChainsQ(r).SelectIDs()...).
		FilterBySrcChain(req.FilterSrcChain).
		FilterByCreator(req.FilterCreator).
		FilterByState(req.FilterState).
		FilterExpired(req.FilterExpired).
		Page(&req.CursorPageParams).
		Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get match orders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var orders []resources.Order
	chains := make([]resources.Chain, 0, 2*len(matches))
	matchesRes := make([]resources.Match, 0, len(matches))

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
	}

	if req.IncludeOriginOrder{
		ids:=make([]int64, len(matches))
		for i,m:=range matches{
			ids[i]=m.OrderID
		}

		dbo,err:=OrdersQ(r).FilterByOrderID(ids...).Select()
		if err!=nil{
			Log(r).WithError(err).Error("failed to get match orders")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		for _,o:=range dbo{
			src := ChainsQ(r).FilterByChainID(o.SrcChain).Get()
			dest := ChainsQ(r).FilterByChainID(o.DestChain).Get()
			orders=append(orders, responses.ToOrderResource(o, src.Key, dest.Key))
		}
	}

	var last string
	if len(matches) > 0 {
		last = strconv.FormatInt(matches[len(matches)-1].ID, 10)
	}

	resp := responses.NewMatchList(matchesRes, orders, chains)
	resp.Links = req.GetCursorLinks(r, last)
	ape.Render(w, resp)
}
