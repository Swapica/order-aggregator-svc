package handlers

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func SetBlock(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewSetBlock(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if err = BlockQ(r).Set(request.Number, request.Chain); err != nil {
		Log(r).WithFields(logan.F{"number": request.Number, "chain": request.Chain}).
			WithError(err).Error("failed to set last block number")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
