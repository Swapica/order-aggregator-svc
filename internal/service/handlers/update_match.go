package handlers

import (
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"github.com/Swapica/order-aggregator-svc/internal/ws"
	"github.com/Swapica/order-aggregator-svc/resources"
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

	match, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if match == nil {
		log.Warn("match order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	if match.State == data.StateBadToken {
		log.Info("match order was hidden due to invalid token_to_buy or token_to_sell, its state won't be updated")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	newState := request.Body.Data.Attributes.State
	if err = q.Update(newState); err != nil {
		log.WithError(err).Error("failed to update match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)

	match.State = newState
	matchResponse := responses.ToMatchResource(
		*match,
		resources.NewKeyInt64(match.SrcChain, "chain"),
		resources.NewKeyInt64(match.OriginOrder, "chain"),
	)
	err = WebSocket(r).BroadcastToClients(ws.UpdateMatch, matchResponse)
	if err != nil {
		log.WithError(err).Debug("failed to broadcast update match order to websocket")
	}
}
