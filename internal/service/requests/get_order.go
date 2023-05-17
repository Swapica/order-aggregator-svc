package requests

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"strconv"
)

type GetOrderRequest struct {
	OrderId int64
}

func NewGetOrder(r *http.Request) (GetOrderRequest, error) {
	var request GetOrderRequest

	orderId := chi.URLParam(r, "id")
	orderIdInt, err := strconv.ParseInt(orderId, 10, 64)
	if err != nil {
		return GetOrderRequest{}, errors.Wrap(err, "failed to parse match id to int")
	}
	request.OrderId = orderIdInt

	return request, nil
}
