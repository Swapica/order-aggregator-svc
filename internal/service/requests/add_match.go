package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type AddMatchRequest resources.MatchResponse

func NewAddMatchRequest(r *http.Request) (*AddMatchRequest, error) {
	var dst AddMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}

	return &dst, dst.validate()
}

func (r *AddMatchRequest) validate() error {
	a := r.Data.Attributes
	return val.Errors{
		"data/id":                       val.Validate(r.Data.ID, val.Required, val.Match(uint256Regexp)),
		"data/type":                     val.Validate(r.Data.Type, val.Required, val.In(resources.MATCH_ORDER)),
		"data/attributes/srcChain":      val.Validate(a.SrcChain, val.Required),
		"data/attributes/originOrderId": val.Validate(a.OriginOrderId, val.Required),
		"data/attributes/account":       val.Validate(a.Account, val.Required, val.Match(addressRegexp)),
		"data/attributes/tokensToSell":  val.Validate(a.TokensToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/amountToSell":  val.Validate(a.AmountToSell.String(), val.Required, val.Match(uint256Regexp)),
		"data/attributes/originChain":   val.Validate(a.OriginChain, val.Required),
		"data/attributes/state":         val.Validate(a.State, val.Required, val.Min(uint8(1))),
	}.Filter()
}

func (r *AddMatchRequest) DBModel() data.Match {
	a := r.Data.Attributes
	return data.Match{
		ID:            r.Data.ID,
		SrcChain:      a.SrcChain,
		OriginOrderId: a.OriginOrderId.String(),
		Account:       a.Account,
		TokensToSell:  a.TokensToSell,
		AmountToSell:  a.AmountToSell.String(),
		OriginChain:   a.OriginChain.String(),
		State:         a.State,
	}
}
