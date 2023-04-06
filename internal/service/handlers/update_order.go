package handlers

import (
	"database/sql"
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

func UpdateOrder(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewUpdateOrder(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	q := OrdersQ(r).FilterByOrderID(req.OrderID).FilterBySrcChain(&req.Chain)
	log := Log(r).WithFields(logan.F{"order_id": req.OrderID, "src_chain": req.Chain})

	order, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if order == nil {
		log.Warn("order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	if order.State == data.StateBadToken {
		log.Info("order was hidden due to invalid token_to_buy or token_to_sell, its status won't be updated")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var matchFK *int64
	matchId := req.Body.Data.Attributes.MatchId
	if matchId != nil {
		log = log.WithField("match_id", *matchId)
		by, err := MatchOrdersQ(r).FilterByMatchID(*matchId).FilterBySrcChain(&order.DestChain).Get()
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

	if matchId != nil {
		order.MatchID = sql.NullInt64{Int64: *matchId, Valid: true}
	}
	if matchFK != nil {
		order.ExecutedByMatch = sql.NullInt64{Int64: *matchFK, Valid: true}
	}
	if a.MatchSwapica != nil {
		order.MatchSwapica = sql.NullString{
			String: *a.MatchSwapica,
			Valid:  true,
		}
	}

	orderResponse := responses.ToOrderResource(
		*order,
		resources.NewKeyInt64(order.SrcChain, "chain"),
		resources.NewKeyInt64(order.DestChain, "chain"),
	)
	err = WebSocket(r).BroadcastToClients(ws.UpdateOrder, orderResponse)
	if err != nil {
		log.WithError(err).Debug("failed to broadcast update order to websocket")
	}
}
