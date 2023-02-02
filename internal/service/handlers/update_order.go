package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateOrder(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := OrdersQ(r).FilterByOrderID(request.OrderID).FilterByChain(&request.Chain)
	log := Log(r).WithFields(logan.F{"order_id": request.OrderID, "src_chain": request.Chain})

	exists, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if exists == nil {
		log.Debug("order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	a := request.Body.Data.Attributes
	if err = q.Update(a.State, request.ExecutedBy, a.MatchSwapica); err != nil {
		log.WithError(err).Error("failed to insert order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
