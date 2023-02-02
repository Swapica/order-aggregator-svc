package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type ListMatches struct {
	page.CursorParams
	FilterChain   *int64  `filter:"chain"`
	FilterState   *uint8  `filter:"state"`
	FilterAccount *string `filter:"account"`
	FilterExpired *bool   `filter:"expired"`
}

func NewListMatches(r *http.Request) (*ListMatches, error) {
	var dst ListMatches
	if err := urlval.Decode(r.URL.Query(), &dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request URL params")
	}

	return &dst, dst.Validate()
}
