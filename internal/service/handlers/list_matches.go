package handlers

import (
	"net/http"
	"strconv"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func ListMatches(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewListMatches(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	matches, err := MatchOrdersQ(r).
		FilterByChain(req.FilterChain).
		FilterByCreator(req.FilterCreator).
		FilterByState(req.FilterState).
		FilterExpired(req.FilterExpired).
		Page(&req.CursorPageParams).
		Select()
	if err != nil {
		Log(r).WithError(err).Error("failed to get match orders")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var last string
	if len(matches) > 0 {
		last = strconv.FormatInt(matches[len(matches)-1].ID, 10)
	}

	resp := responses.NewMatchList(matches)
	resp.Links = req.GetCursorLinks(r, last)
	ape.Render(w, resp)
}
