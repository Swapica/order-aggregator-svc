package handlers

import (
	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
)

func GetOrder(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetOrder(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	order, err := OrdersQ(r).New().FilterByOrderID(request.OrderId).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get order by id")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if order == nil {
		Log(r).WithFields(logan.F{
			"order_id": request.OrderId,
		}).Error("order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, order)
}
