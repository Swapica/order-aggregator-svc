package handlers

import (
	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
	"net/http"
)

func GetMatch(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetMatch(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	match, err := MatchOrdersQ(r).New().FilterByMatchID(request.MatchId).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get match by id")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if match == nil {
		Log(r).WithFields(logan.F{
			"match_id": request.MatchId,
		}).Error("match not found")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, match)
}
