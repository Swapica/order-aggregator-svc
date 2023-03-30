package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/internal/service/helpers"
	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"github.com/Swapica/order-aggregator-svc/internal/ws"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func AddMatch(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewAddMatch(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	attr := req.Data.Attributes
	q := MatchOrdersQ(r).FilterByMatchID(attr.MatchId).FilterBySrcChain(&attr.SrcChainId)
	log := Log(r).WithFields(logan.F{
		"match_id": attr.MatchId, "src_chain": attr.SrcChainId,
		"origin_order_id": attr.OriginOrderId, "origin_chain_id": attr.OriginChainId})

	conflict, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if conflict != nil {
		log.Info("match order already exists")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	originOrder, err := OrdersQ(r).FilterByOrderID(attr.OriginOrderId).FilterBySrcChain(&attr.OriginChainId).Get()
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

	srcChain := ChainsQ(r).FilterByChainID(attr.SrcChainId).Get()
	if srcChain == nil {
		log.Warn("src_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	originChain := ChainsQ(r).FilterByChainID(attr.OriginChainId).Get()
	if originChain == nil {
		log.Warn("origin_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	sellToken, err := helpers.GetOrAddToken(TokensQ(r), attr.TokenToSell, *srcChain)
	// token_to_sell == origin_order.token_to_buy, therefore assertion of order state covers the check for a bad token
	if err != nil && !helpers.IsBadTokenErr(err) {
		log.WithError(err).Error("failed to get or add token to sell")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	match := req.DBModel(originOrder.ID, sellToken.ID)
	if originOrder.State == data.StateBadToken {
		log.Info("origin_order has invalid token_to_buy or token_to_sell, the match will be hidden")
		match.State = data.StateBadToken
	}

	match, err = q.Insert(match)
	if err != nil {
		log.WithError(err).Error("failed to add match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusCreated)
	ape.Render(w, responses.NewMatch(match, srcChain.Key, originChain.Key))

	matchResponse := responses.ToMatchResource(match, srcChain.Key, originChain.Key)
	err = WebSocket(r).BroadcastToClients(ws.AddMatch, matchResponse)
	if err != nil {
		log.WithError(err).Debug("failed to broadcast match order to websocket")
	}
}
