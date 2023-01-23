package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateOrderRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := OrdersQ(r)
	id, chain := request.Data.ID, request.Data.Attributes.SrcChain
	log := Log(r).WithFields(logan.F{"order_id": id, "src_chain": chain})

	exists, err := q.Get(id, chain)
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

	a := request.Data.Attributes
	if err = q.Update(id, chain, a.State, a.ExecutedBy, a.MatchSwapica); err != nil {
		log.WithError(err).Error("failed to insert order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
