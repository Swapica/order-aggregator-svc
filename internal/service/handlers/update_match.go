package handlers

import (
	"fmt"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/internal/service/notifications"
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

	if err = q.Update(request.Body.Data.Attributes.State); err != nil {
		log.WithError(err).Error("failed to update match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	// TODO check chains everywhere
	originOrder, err := OrdersQ(r).FilterByOrderID(match.OriginOrder).FilterBySrcChain(&match.SrcChain).Get()
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

	srcChain := ChainsQ(r).FilterByChainID(match.SrcChain).Get()
	if srcChain == nil {
		log.Warn("src_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	originChain := ChainsQ(r).FilterByChainID(match.OrderChain).Get()
	if originChain == nil {
		log.Warn("origin_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	// FIXME chain? maybe just to Ethereum?
	pushCli := notifications.NewNotificationsClient(Notifications(r), match.SrcChain)

	if err := pushCli.NotifyUser(
		fmt.Sprintf("Match for the %s/%s order has been updated",
			t[0].Symbol, t[1].Symbol),
		fmt.Sprintf("Sell amount: %s.\nBuy amount: %s.\nSource chain: %s.\nDestination chain: %s.\nState: %s.\n",
			originOrder.SellAmount, originOrder.BuyAmount,
			srcChain.Attributes.Name, originChain.Attributes.Name,
			data.StateToString(match.State)),
		originOrder.Creator,
	); err != nil {
		log.WithError(err).Error("failed to notify user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
