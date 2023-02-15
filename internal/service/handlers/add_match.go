package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func AddMatch(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewAddMatch(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	match := request.DBModel()
	q := MatchOrdersQ(r).FilterByMatchID(match.MatchID).FilterByChain(&match.SrcChain)
	log := Log(r).WithFields(logan.F{
		"match_id": match.MatchID, "src_chain": match.SrcChain, "order_id": match.OrderID, "order_chain": match.OrderChain})

	conflict, err := q.Get()
	if err != nil {
		log.WithError(err).Error("failed to get match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if conflict != nil {
		log.Debug("match order already exists")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	origin, err := OrdersQ(r).FilterByOrderID(match.OrderID).FilterByChain(&match.OrderChain).Get()
	if err != nil {
		log.WithError(err).Error("failed to get origin order")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if origin == nil {
		log.Debug("origin order not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	newMatch, err := q.Insert(match)
	if err != nil {
		log.WithError(err).Error("failed to add match order")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusCreated)
	ape.Render(w, responses.NewMatch(newMatch))
}
