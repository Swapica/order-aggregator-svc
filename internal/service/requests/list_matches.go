package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type ListMatchesRequest struct {
	Chain string
	page.CursorParams
	FilterState   *string `filter:"state"`
	FilterAccount *string `filter:"account"`
	FilterExpired *bool   `filter:"expired"`
}

func NewListMatchesRequest(r *http.Request) (*ListMatchesRequest, error) {
	dst := ListMatchesRequest{Chain: chi.URLParam(r, "chain")}
	if err := validateChain(dst.Chain); err != nil {
		return nil, err
	}

	if err := urlval.Decode(r.URL.Query(), &dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request URL params")
	}

	return &dst, dst.validate()
}

func (r *ListMatchesRequest) validate() error {
	if err := r.CursorParams.Validate(); err != nil {
		return err
	}
	return val.Errors{
		"filter[account]": val.Validate(r.FilterAccount, val.Match(addressRegexp)),
		"filter[state]":   validateState(r.FilterState),
	}.Filter()
}
