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
		FilterByState(req.FilterState)

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

	chains := make([]resources.Chain, 0, 2*len(orders))
	ordersRes := make([]resources.Order, 0, len(orders))
	for _, o := range orders {
		// must not be nil because of FilterBySupportedChains
		src := ChainsQ(r).FilterByChainID(o.SrcChain).Get()
		dest := ChainsQ(r).FilterByChainID(o.DestChain).Get()
		ordersRes = append(ordersRes, responses.ToOrderResource(o, src.Key, dest.Key))

		if req.IncludeSrcChain {
			chains = append(chains, *src)
		}
		if req.IncludeDestChain {
			chains = append(chains, *dest)
		}
	}

	resp := responses.NewOrderList(ordersRes, chains, count)
	resp.Links = req.Params.GetLinks(r)
	ape.Render(w, resp)
}
