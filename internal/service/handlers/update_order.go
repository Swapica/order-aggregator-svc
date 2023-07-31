package handlers

import (
	"database/sql"
	"fmt"
	"math/big"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/internal/service/helpers"
	"github.com/Swapica/order-aggregator-svc/internal/service/notifications"
	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"github.com/Swapica/order-aggregator-svc/internal/ws"
	"github.com/Swapica/order-aggregator-svc/resources"
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
	var match *data.Match
	matchId := req.Body.Data.Attributes.MatchId
	if matchId != nil {
		log = log.WithField("match_id", *matchId)
		match, err = MatchOrdersQ(r).FilterByMatchID(*matchId).FilterBySrcChain(&order.DestChain).Get()
		if err != nil {
			log.WithError(err).Error("failed to get match order that executed the order")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if match == nil {
			log.Warn("match order that executed the order not found")
			ape.RenderErr(w, problems.NotFound())
			return
		}
		matchFK = &match.ID
	}

	a := req.Body.Data.Attributes
	if err = q.Update(a.State, matchFK, matchId, a.MatchSwapica); err != nil {
		log.WithError(err).Error("failed to update order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if matchId != nil {
		t, err := TokensQ(r).New().FilterByID(order.SellToken, order.BuyToken).Select()
		if err != nil {
			log.WithError(err).Error("failed to get token from database")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		if len(t) < 2 {
			log.Error("token(s) not found")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		matchSrcChain := ChainsQ(r).FilterByChainID(match.SrcChain).Get()
		if matchSrcChain == nil {
			log.Warn("src_chain is not supported by swapica-svc")
			ape.RenderErr(w, problems.NotFound())
			return
		}
		matchDestChain := ChainsQ(r).FilterByChainID(match.OrderChain).Get()
		if matchDestChain == nil {
			log.Warn("origin_chain is not supported by swapica-svc")
			ape.RenderErr(w, problems.NotFound())
			return
		}

		pushCli := notifications.NewNotificationsClient(Notifications(r), 1)

		sellAmountI, _ := new(big.Int).SetString(order.SellAmount, 10)
		buyAmountI, _ := new(big.Int).SetString(order.BuyAmount, 10)

		if err := pushCli.NotifyUser(
			fmt.Sprintf("The %s/%s order you matched has been updated",
				t[0].Symbol, t[1].Symbol),
			fmt.Sprintf("Order sell amount: %s.\nOrder buy amount: %s.\nMatch source chain: %s.\nMatch destination chain: %s.\nOrder state: %s.\n",
				helpers.ConvertAmount(sellAmountI, t[0].Decimals).String(),
				helpers.ConvertAmount(buyAmountI, t[1].Decimals).String(),
				matchSrcChain.Attributes.Name, matchDestChain.Attributes.Name,
				data.StateToString(a.State)),
			match.Creator,
		); err != nil {
			log.WithError(err).Error("failed to notify user")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		if order.UseRelayer && a.State == 4 {
			if err := pushCli.NotifyUser(
				fmt.Sprintf("Your %s/%s order has been executed",
					t[0].Symbol, t[1].Symbol),
				fmt.Sprintf("Order sell amount: %s.\nOrder buy amount: %s.\nMatch source chain: %s.\nMatch destination chain: %s.\nOrder state: %s.\n",
					helpers.ConvertAmount(sellAmountI, t[0].Decimals).String(),
					helpers.ConvertAmount(buyAmountI, t[1].Decimals).String(),
					matchSrcChain.Attributes.Name, matchDestChain.Attributes.Name,
					data.StateToString(a.State)),
				order.Creator,
			); err != nil {
				log.WithError(err).Error("failed to notify user")
				ape.RenderErr(w, problems.InternalError())
				return
			}
		}
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
