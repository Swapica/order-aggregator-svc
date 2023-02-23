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

	var matchFK *int64
	matchId := req.Body.Data.Attributes.MatchId
	if matchId != nil {
		log = log.WithField("match_id", *matchId)
		by, err := MatchOrdersQ(r).FilterByMatchID(*matchId).FilterBySrcChain(&exists.DestChain).Get()
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
		matchFK = &by.ID
	}

	a := req.Body.Data.Attributes
	if err = q.Update(a.State, matchFK, matchId, a.MatchSwapica); err != nil {
		log.WithError(err).Error("failed to update order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
