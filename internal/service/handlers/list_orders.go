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

func ListOrders(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewListOrders(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	orders, err := OrdersQ(r).
		FilterByChain(req.FilterChain).
		FilterByDestinationChain(req.FilterDestChain).
		FilterByTokenToBuy(req.FilterBuyToken).
		FilterByTokenToSell(req.FilterSellToken).
		FilterByCreator(req.FilterCreator).
		FilterByState(req.FilterState).
		Page(&req.CursorPageParams).
		Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get orders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var chains []resources.Chain
	chainIDs := make([]int64, 0, 2*len(orders))
	if req.IncludeSrcChain || req.IncludeDestChain {
		for _, o := range orders {
			if req.IncludeSrcChain {
				chainIDs = append(chainIDs, o.SrcChain)
			}
			if req.IncludeDestChain {
				chainIDs = append(chainIDs, o.DestChain)
			}
		}
		chains = ChainsQ(r).FilterByChainID(chainIDs...).Select()
	}

	var last string
	if len(orders) > 0 {
		last = strconv.FormatInt(orders[len(orders)-1].ID, 10)
	}

	resp := responses.NewOrderList(orders, chains)
	resp.Links = req.GetCursorLinks(r, last)
	ape.Render(w, resp)
}
