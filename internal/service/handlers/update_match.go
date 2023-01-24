package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func UpdateMatch(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateMatchRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := MatchOrdersQ(r)
	id := request.Body.Data.ID
	log := Log(r).WithFields(logan.F{"match_id": id, "src_chain": request.Chain})

	exists, err := q.Get(id, request.Chain)
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

	if err = q.Update(id, request.Chain, request.Body.Data.Attributes.State); err != nil {
		log.WithError(err).Error("failed to insert match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
