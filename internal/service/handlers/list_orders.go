package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"github.com/Swapica/order-aggregator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ListOrders(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewListOrders(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := OrdersQ(r).
		FilterBySupportedChains(ChainsQ(r).SelectIDs()...).
		FilterBySrcChain(req.FilterSrcChain).
		FilterByDestChain(req.FilterDestChain).
		FilterByTokenToBuy(req.FilterBuyToken).
		FilterByTokenToSell(req.FilterSellToken).
		FilterByCreator(req.FilterCreator).
		FilterByState(req.FilterState).
		FilterByUseRelayer(req.FilterUseRelayer)

	orders, err := q.Page(&req.OffsetPageParams).Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get orders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	count, err := q.Count()
	if err != nil {
		Log(r).WithError(err).Error("failed to count orders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ordersRes := make([]resources.Order, 0, len(orders))
	tokenIDs := make([]int64, 0, 3*len(orders))
	included := make([]resources.Resource, 0, 4*len(orders))

	for _, o := range orders {
		// must not be nil because of FilterBySupportedChains
		src := ChainsQ(r).FilterByChainID(o.SrcChain).Get()
		dest := ChainsQ(r).FilterByChainID(o.DestChain).Get()
		ordersRes = append(ordersRes, responses.ToOrderResource(o, src.Key, dest.Key))

		if req.IncludeSrcChain {
			included = append(included, src)
		}
		if req.IncludeDestChain {
			included = append(included, dest)
		}
		if req.IncludeBuyToken {
			tokenIDs = append(tokenIDs, o.BuyToken)
		}
		if req.IncludeSellToken {
			tokenIDs = append(tokenIDs, o.SellToken)
		}
	}

	if req.IncludeBuyToken || req.IncludeSellToken {
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

	resp := responses.NewOrderList(ordersRes, included, count)
	resp.Links = req.Params.GetLinks(r)
	ape.Render(w, resp)
}
