package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/internal/service/responses"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBlock(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetBlock(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log := Log(r).WithField("chain", request.Chain)
	num, err := BlockQ(r).Get(request.Chain)
	if err != nil {
		log.WithError(err).Error("failed to get last block number")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if num == nil {
		log.Debug("last block is not set")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, responses.NewBlock(*num))
}
