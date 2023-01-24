package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type UpdateMatchRequest struct {
	Body  resources.UpdateMatchRequest
	Chain string
}

func NewUpdateMatchRequest(r *http.Request) (*UpdateMatchRequest, error) {
	dst := UpdateMatchRequest{Chain: chi.URLParam(r, "chain")}
	if err := json.NewDecoder(r.Body).Decode(&dst.Body); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}

	a := dst.Body.Data.Attributes
	return &dst, val.Errors{
		"{chain}":               val.Validate(dst.Chain, val.Required),
		"data/id":               val.Validate(dst.Body.Data.ID, val.Required, val.Match(uint256Regexp)),
		"data/type":             val.Validate(dst.Body.Data.Type, val.Required, val.In(resources.MATCH_ORDER)),
		"data/attributes/state": val.Validate(a.State, val.Required, val.Min(uint8(1))),
	}.Filter()
}
