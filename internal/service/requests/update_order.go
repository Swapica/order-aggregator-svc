package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type UpdateOrderRequest struct {
	Body  resources.UpdateOrderRequest
	Chain string
}

func NewUpdateOrderRequest(r *http.Request) (*UpdateOrderRequest, error) {
	dst := UpdateOrderRequest{Chain: chi.URLParam(r, "chain")}
	if err := json.NewDecoder(r.Body).Decode(&dst.Body); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}

	a := dst.Body.Data.Attributes
	return &dst, val.Errors{
		"{chain}":                      validateUint(dst.Chain, bigintBitSize),
		"data/id":                      validateUint(dst.Chain, bigintBitSize),
		"data/type":                    val.Validate(dst.Body.Data.Type, val.Required, val.In(resources.ORDER)),
		"data/attributes/state":        val.Validate(a.State, val.Required, val.Min(uint8(1))),
		"data/attributes/matchSwapica": val.Validate(a.MatchSwapica, val.NilOrNotEmpty, val.Match(addressRegexp)),
	}.Filter()
}
