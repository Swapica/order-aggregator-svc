package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewUpdateOrder(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := OrdersQ(r).FilterByOrderID(req.OrderID).FilterBySrcChain(&req.Chain)
	log := Log(r).WithFields(logan.F{"order_id": req.OrderID, "src_chain": req.Chain})

	exists, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if exists == nil {
		log.Warn("order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if req.MatchID != nil {
		log = log.WithField("match_id", *req.MatchID)
		by, err := MatchOrdersQ(r).FilterByMatchID(*req.MatchID).FilterBySrcChain(&exists.DestChain).Get()
		if err != nil {
			log.WithError(err).Error("failed to get match order that executed the order")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if by == nil {
			log.Warn("match order that executed the order not found")
			ape.RenderErr(w, problems.NotFound())
			return
		}
	}

	a := req.Body.Data.Attributes
	if err = q.Update(a.State, req.MatchID, a.MatchSwapica); err != nil {
		log.WithError(err).Error("failed to insert order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
