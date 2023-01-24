package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AddMatchRequest struct {
	Body  resources.MatchResponse
	Chain string
}

func NewAddMatchRequest(r *http.Request) (*AddMatchRequest, error) {
	dst := AddMatchRequest{Chain: chi.URLParam(r, "chain")}
	if err := json.NewDecoder(r.Body).Decode(&dst.Body); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}

	return &dst, dst.validate()
}

func (r *AddMatchRequest) validate() error {
	a := r.Body.Data.Attributes
	origin := r.Body.Data.Relationships.OriginOrder.Data
	if origin == nil {
		return val.Errors{"data/relationships/originOrder/data": val.Validate(origin, val.NotNil)}
	}

	return val.Errors{
		"{chain}":                                  val.Validate(r.Chain, val.Required),
		"data/id":                                  val.Validate(r.Body.Data.ID, val.Required, val.Match(uint256Regexp)),
		"data/type":                                val.Validate(r.Body.Data.Type, val.Required, val.In(resources.MATCH_ORDER)),
		"data/attributes/account":                  val.Validate(a.Account, val.Required, val.Match(addressRegexp)),
		"data/attributes/tokenToSell":              val.Validate(a.TokenToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/amountToSell":             val.Validate(a.AmountToSell.String(), val.Required, val.Match(uint256Regexp)),
		"data/attributes/originChain":              val.Validate(a.OriginChain, val.Required),
		"data/attributes/state":                    val.Validate(a.State, val.Required, val.Min(uint8(1))),
		"data/relationships/originOrder/data/id":   val.Validate(origin.ID, val.Required, val.Match(uint256Regexp)),
		"data/relationships/originOrder/data/type": val.Validate(origin.Type, val.Required, val.In(resources.ORDER)),
	}.Filter()
}

func (r *AddMatchRequest) DBModel() data.Match {
	a := r.Body.Data.Attributes
	return data.Match{
		ID:            r.Body.Data.ID,
		OriginOrderId: r.Body.Data.Relationships.OriginOrder.Data.ID,
		SrcChain:      r.Chain,
		Account:       a.Account,
		TokenToSell:   a.TokenToSell,
		AmountToSell:  a.AmountToSell.String(),
		OriginChain:   a.OriginChain.String(),
		State:         a.State,
	}
}
