package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func NewUpdateOrderRequest(r *http.Request) (*resources.UpdateOrderRequest, error) {
	var dst resources.UpdateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}

	a := dst.Data.Attributes
	return &dst, val.Errors{
		"data/id":                      val.Validate(dst.Data.ID, val.Required, val.Match(uint256Regexp)),
		"data/type":                    val.Validate(dst.Data.Type, val.Required, val.In(resources.ORDER)),
		"data/attributes/srcChain":     val.Validate(a.SrcChain, val.Required),
		"data/attributes/state":        val.Validate(a.State, val.Required, val.Min(uint8(1))),
		"data/attributes/matchSwapica": val.Validate(a.MatchSwapica, val.NilOrNotEmpty, val.Match(addressRegexp)),
	}.Filter()
}
