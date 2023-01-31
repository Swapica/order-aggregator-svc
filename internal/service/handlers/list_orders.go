package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ListOrders(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewListOrders(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	orders, err := OrdersQ(r).FilterByChain(req.Chain).
		FilterByTokenToBuy(req.FilterBuyToken).
		FilterByTokenToSell(req.FilterSellToken).
		FilterByAccount(req.FilterAccount).
		FilterByState(req.FilterState).
		Page(&req.CursorPageParams).
		Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get orders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var last string
	if len(orders) > 0 {
		last = orders[len(orders)-1].ID
	}

	resp := responses.NewOrderList(orders)
	resp.Links = req.GetCursorLinks(r, last)
	ape.Render(w, resp)
}
