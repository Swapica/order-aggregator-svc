package handlers

import (
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

	originOrder, err := OrdersQ(r).FilterByOrderID(match.OrderID).FilterBySrcChain(&match.OrderChain).Get()
	if err != nil {
		log.WithError(err).Error("failed to get origin order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if originOrder == nil {
		log.Warn("origin order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	t, err := TokensQ(r).New().FilterByID(originOrder.SellToken, originOrder.BuyToken).Select()
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

	orderSrcChain := ChainsQ(r).FilterByChainID(match.OrderChain).Get()
	if orderSrcChain == nil {
		log.Warn("src_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	orderDestChain := ChainsQ(r).FilterByChainID(match.SrcChain).Get()
	if orderDestChain == nil {
		log.Warn("origin_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	pushCli := notifications.NewNotificationsClient(Notifications(r), 1)

	sellAmountI, _ := new(big.Int).SetString(originOrder.SellAmount, 10)
	buyAmountI, _ := new(big.Int).SetString(originOrder.BuyAmount, 10)

	if err := pushCli.NotifyUser(
		fmt.Sprintf("Match for the %s/%s order has been updated",
			t[0].Symbol, t[1].Symbol),
		fmt.Sprintf("Order sell amount: %s.\nOrder buy amount: %s.\nOrder source chain: %s.\nOrder destination chain: %s.\nMatch state: %s.\n",
			helpers.ConvertAmount(sellAmountI, t[0].Decimals).String(),
			helpers.ConvertAmount(buyAmountI, t[1].Decimals).String(),
			orderSrcChain.Attributes.Name, orderDestChain.Attributes.Name,
			data.StateToString(newState)),
		originOrder.Creator,
	); err != nil {
		log.WithError(err).Error("failed to notify user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if match.UseRelayer && newState == 4 {
		if err := pushCli.NotifyUser(
			fmt.Sprintf("Your match for the %s/%s order has been executed",
				t[0].Symbol, t[1].Symbol),
			fmt.Sprintf("Order sell amount: %s.\nOrder buy amount: %s.\nOrder source chain: %s.\nOrder destination chain: %s.\nMatch state: %s.\n",
				helpers.ConvertAmount(sellAmountI, t[0].Decimals).String(),
				helpers.ConvertAmount(buyAmountI, t[1].Decimals).String(),
				orderSrcChain.Attributes.Name, orderDestChain.Attributes.Name,
				data.StateToString(newState)),
			match.Creator,
		); err != nil {
			log.WithError(err).Error("failed to notify user")
			ape.RenderErr(w, problems.InternalError())
			return
		}
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
