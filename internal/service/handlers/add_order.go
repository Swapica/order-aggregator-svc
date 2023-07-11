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

func AddOrder(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewAddOrder(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	attr := req.Data.Attributes
	q := OrdersQ(r).FilterByOrderID(attr.OrderId).FilterBySrcChain(&attr.SrcChainId)
	log := Log(r).WithFields(logan.F{"order_id": attr.OrderId, "src_chain": attr.SrcChainId, "dest_chain": attr.DestChainId})
	badTokens := false

	conflict, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if conflict != nil {
		log.Info("order already exists")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	srcChain := ChainsQ(r).FilterByChainID(attr.SrcChainId).Get()
	if srcChain == nil {
		log.Warn("src_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}
	destChain := ChainsQ(r).FilterByChainID(attr.DestChainId).Get()
	if destChain == nil {
		log.Warn("dest_chain is not supported by swapica-svc")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	sellToken, err := helpers.GetOrAddToken(TokensQ(r), attr.TokenToSell, *srcChain)
	if err != nil {
		if !helpers.IsBadTokenErr(err) {
			log.WithError(err).Error("failed to get or add token to sell")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		log.WithError(err).Info("found bad token_to_sell, the order will be hidden")
		badTokens = true
	}

	buyToken, err := helpers.GetOrAddToken(TokensQ(r), attr.TokenToBuy, *destChain)
	if err != nil {
		if !helpers.IsBadTokenErr(err) {
			log.WithError(err).Error("failed to get or add token to sell")
			ape.RenderErr(w, problems.InternalError())
			return
		}
		log.WithError(err).Info("found invalid token_to_buy, the order will be hidden")
		badTokens = true
	}

	order := req.DBModel(sellToken.ID, buyToken.ID)
	if badTokens {
		order.State = data.StateBadToken
	}

	order, err = q.Insert(req.DBModel(sellToken.ID, buyToken.ID))
	if err != nil {
		log.WithError(err).Error("failed to add order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusCreated)
	ape.Render(w, responses.NewOrder(order, srcChain.Key, destChain.Key))

	orderResponse := responses.ToOrderResource(order, srcChain.Key, destChain.Key)
	err = WebSocket(r).BroadcastToClients(ws.AddOrder, orderResponse)
	if err != nil {
		log.WithError(err).Debug("failed to broadcast order to websocket")
	}
}
