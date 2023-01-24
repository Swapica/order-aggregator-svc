package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func AddMatch(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddMatchRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := MatchOrdersQ(r)
	match := request.DBModel()
	log := Log(r).WithFields(logan.F{
		"match_id": match.ID, "origin_order_id": match.OriginOrderId, "src_chain": match.SrcChain})

	conflict, err := q.Get(match.ID, match.SrcChain)
	if err != nil {
		log.WithError(err).Error("failed to get match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if conflict != nil {
		log.Debug("match order already exists")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	origin, err := OrdersQ(r).Get(match.OriginOrderId, match.SrcChain)
	if err != nil {
		log.WithError(err).Error("failed to get origin order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if origin == nil {
		log.Debug("origin order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if err = q.Insert(match); err != nil {
		log.WithError(err).Error("failed to insert match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
