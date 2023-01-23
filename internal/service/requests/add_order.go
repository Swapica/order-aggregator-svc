package requests

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var uint256Regexp = regexp.MustCompile(`^[0-9]{1,78}$`)
var addressRegexp = regexp.MustCompile("^0x[0-9A-Fa-f]{40}$")

type AddOrderRequest resources.OrderResponse

func NewAddOrderRequest(r *http.Request) (*AddOrderRequest, error) {
	var dst AddOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}

	return &dst, dst.validate()
}

func (r *AddOrderRequest) validate() error {
	a := r.Data.Attributes
	return val.Errors{
		"data/id":                      val.Validate(r.Data.ID, val.Required, val.Match(uint256Regexp)),
		"data/type":                    val.Validate(r.Data.Type, val.Required, val.In(resources.ORDER)),
		"data/attributes/srcChain":     val.Validate(a.SrcChain, val.Required),
		"data/attributes/account":      val.Validate(a.Account, val.Required, val.Match(addressRegexp)),
		"data/attributes/tokensToSell": val.Validate(a.TokensToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/tokensToBuy":  val.Validate(a.TokensToBuy, val.Required, val.Match(addressRegexp)),
		"data/attributes/amountToSell": val.Validate(a.AmountToSell.String(), val.Required, val.Match(uint256Regexp)),
		"data/attributes/amountToBuy":  val.Validate(a.AmountToBuy.String(), val.Required, val.Match(uint256Regexp)),
		"data/attributes/destChain":    val.Validate(a.DestChain, val.Required),
		"data/attributes/state":        val.Validate(a.State, val.Required, val.Min(uint8(1))),
		"data/attributes/matchSwapica": val.Validate(a.MatchSwapica, val.NilOrNotEmpty, val.Match(addressRegexp)),
	}.Filter()
}

func (r *AddOrderRequest) DBModel() data.Order {
	a := r.Data.Attributes
	order := data.Order{
		ID:           r.Data.ID,
		SrcChain:     a.SrcChain,
		Account:      a.Account,
		TokensToSell: a.TokensToSell,
		TokensToBuy:  a.TokensToBuy,
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
