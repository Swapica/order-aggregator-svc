package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func AddOrder(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddOrder(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	order := request.DBModel()
	q := OrdersQ(r).FilterByOrderID(order.OrderID).FilterByChain(&order.SrcChain)
	log := Log(r).WithFields(logan.F{"order_id": order.OrderID, "src_chain": order.SrcChain})

	conflict, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if conflict != nil {
		log.Debug("order already exists")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	if err = q.Insert(order); err != nil {
		log.WithError(err).Error("failed to insert order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
