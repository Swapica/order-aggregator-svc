package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func UpdateMatch(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateMatch(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := MatchOrdersQ(r).FilterByMatchID(request.MatchID).FilterBySrcChain(&request.Chain)
	log := Log(r).WithFields(logan.F{"match_id": request.MatchID, "src_chain": request.Chain})

	exists, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if exists == nil {
		log.Debug("match order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if err = q.Update(request.Body.Data.Attributes.State); err != nil {
		log.WithError(err).Error("failed to update match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
