package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateMatch struct {
	Body           resources.UpdateMatchRequest
	MatchID, Chain int64
}

func NewUpdateMatch(r *http.Request) (*UpdateMatch, error) {
	var dst UpdateMatch
	if err := json.NewDecoder(r.Body).Decode(&dst.Body); err != nil {
		return nil, toDecodeErr(err, "body")
	}

	var errChain, errMatchID error
	dst.Chain, errChain = parseBigint(chi.URLParam(r, "chain"))
	dst.MatchID, errMatchID = parseBigint(dst.Body.Data.ID)

	return &dst, val.Errors{
		"{chain}":               errChain,
		"data/id":               errMatchID,
		"data/type":             val.Validate(dst.Body.Data.Type, val.Required, val.In(resources.MATCH_ORDER)),
		"data/attributes/state": val.Validate(dst.Body.Data.Attributes.State, val.Required),
	}.Filter()
}
