package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func UpdateMatch(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateMatch(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
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
		log.Warn("match order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	if exists.State == data.StateBadToken {
		log.Info("match order was hidden due to invalid token_to_buy or token_to_sell, its state won't be updated")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err = q.Update(request.Body.Data.Attributes.State); err != nil {
		log.WithError(err).Error("failed to update match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
