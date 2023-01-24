package handlers

import (
	"net/http"
	"strconv"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"github.com/Swapica/order-aggregator-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBlock(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetBlockRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	num, err := BlockQ(r).Get(request.Chain)
	if err != nil {
		Log(r).WithField("chain", request.Chain).
			WithError(err).Error("failed to get last block number")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if num == nil {
		Log(r).Debug("last block is not set")
		ape.RenderErr(w, problems.NotFound())
		return
	}

	ape.Render(w, newBlockResponse(*num))
}

func newBlockResponse(id uint64) resources.BlockResponse {
	return resources.BlockResponse{
		Data: resources.Block{
			Key: resources.Key{
				ID:   strconv.FormatUint(id, 10),
				Type: resources.BLOCK,
			},
		},
	}
}
