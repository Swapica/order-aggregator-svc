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

	match := request.DBModel()
	q := MatchOrdersQ(r).FilterByChain(match.SrcChain)
	log := Log(r).WithFields(logan.F{
		"match_id": match.ID, "src_chain": match.SrcChain, "order_id": match.OrderID, "order_chain": match.OrderChain})

	conflict, err := q.Get(match.ID)
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

	origin, err := OrdersQ(r).FilterByChain(match.OrderChain).Get(match.OrderID)
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
