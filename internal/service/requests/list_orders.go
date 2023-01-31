package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type ListOrders struct {
	Chain string
	page.CursorParams
	FilterBuyToken  *string `filter:"tokenToBuy"`
	FilterSellToken *string `filter:"tokenToSell"`
	FilterAccount   *string `filter:"account"`
	// with *uint8 for values > MaxUint8 it is not decoded correctly
	FilterState *string `filter:"state"`
}

func NewListOrders(r *http.Request) (*ListOrders, error) {
	dst := ListOrders{Chain: chi.URLParam(r, "chain")}
	if err := validateChain(dst.Chain); err != nil {
		return nil, err
	}

	if err := urlval.Decode(r.URL.Query(), &dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request URL params")
	}

	return &dst, dst.validate()
}

func (r *ListOrders) validate() error {
	if err := r.CursorParams.Validate(); err != nil {
		return err
	}
	return val.Errors{
		"filter[tokenToBuy]":  val.Validate(r.FilterBuyToken, val.Match(addressRegexp)),
		"filter[tokenToSell]": val.Validate(r.FilterSellToken, val.Match(addressRegexp)),
		"filter[account]":     val.Validate(r.FilterAccount, val.Match(addressRegexp)),
		"filter[state]":       validateState(r.FilterState),
	}.Filter()
}
