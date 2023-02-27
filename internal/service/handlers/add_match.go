package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/helpers"
	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
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
	if err != nil {
		log.WithError(err).Error("failed to get or add token to sell")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	newMatch, err := q.Insert(req.DBModel(originOrder.ID, sellToken.ID))
	if err != nil {
		log.WithError(err).Error("failed to add match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusCreated)
	ape.Render(w, responses.NewMatch(newMatch, srcChain.Key, originChain.Key))
}
