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
	id, chain := request.Data.ID, request.Data.Attributes.SrcChain
	log := Log(r).WithFields(logan.F{"match_id": id, "src_chain": chain})

	exists, err := q.Get(id, chain)
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

	if err = q.Update(id, chain, request.Data.Attributes.State); err != nil {
		log.WithError(err).Error("failed to insert match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
