package requests

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var uint8Regexp = regexp.MustCompile(`^[0-9]{1,3}$`)
var uint256Regexp = regexp.MustCompile(`^[0-9]{1,78}$`)
var addressRegexp = regexp.MustCompile("^0x[0-9A-Fa-f]{40}$")

type AddOrderRequest struct {
	Body  resources.OrderResponse
	Chain string
}

func NewAddOrderRequest(r *http.Request) (*AddOrderRequest, error) {
	dst := AddOrderRequest{Chain: chi.URLParam(r, "chain")}
	if err := json.NewDecoder(r.Body).Decode(&dst.Body); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}

	return &dst, dst.validate()
}

func (r *AddOrderRequest) validate() error {
	a := r.Body.Data.Attributes
	return val.Errors{
		"{chain}":                      val.Validate(r.Chain, val.Required),
		"data/id":                      val.Validate(r.Body.Data.ID, val.Required, val.Match(uint256Regexp)),
		"data/type":                    val.Validate(r.Body.Data.Type, val.Required, val.In(resources.ORDER)),
		"data/attributes/account":      val.Validate(a.Account, val.Required, val.Match(addressRegexp)),
		"data/attributes/tokenToSell":  val.Validate(a.TokenToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/tokenToBuy":   val.Validate(a.TokenToBuy, val.Required, val.Match(addressRegexp)),
		"data/attributes/amountToSell": val.Validate(a.AmountToSell.String(), val.Required, val.Match(uint256Regexp)),
		"data/attributes/amountToBuy":  val.Validate(a.AmountToBuy.String(), val.Required, val.Match(uint256Regexp)),
		"data/attributes/destChain":    val.Validate(a.DestChain, val.Required),
		"data/attributes/state":        val.Validate(a.State, val.Required, val.Min(uint8(1))),
		"data/attributes/matchSwapica": val.Validate(a.MatchSwapica, val.NilOrNotEmpty, val.Match(addressRegexp)),
	}.Filter()
}

func (r *AddOrderRequest) DBModel() data.Order {
	a := r.Body.Data.Attributes
	order := data.Order{
		ID:           r.Body.Data.ID,
		SrcChain:     r.Chain,
		Account:      a.Account,
		TokenToSell:  a.TokenToSell,
		TokenToBuy:   a.TokenToBuy,
		AmountToSell: a.AmountToSell.String(),
		AmountToBuy:  a.AmountToBuy.String(),
		DestChain:    a.DestChain.String(),
		State:        a.State,
	}

	if a.ExecutedBy != nil {
		order.ExecutedBy = sql.NullString{String: a.ExecutedBy.String(), Valid: true}
	}
	if a.MatchSwapica != nil {
		order.MatchSwapica = sql.NullString{String: *a.MatchSwapica, Valid: true}
	}

	return order
}
