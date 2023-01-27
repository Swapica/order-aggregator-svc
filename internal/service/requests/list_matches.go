package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type ListMatchesRequest struct {
	Chain string
	page.CursorParams
}

func NewListMatchesRequest(r *http.Request) (*ListMatchesRequest, error) {
	dst := ListMatchesRequest{Chain: chi.URLParam(r, "chain")}
	if err := requireChain(dst.Chain); err != nil {
		return nil, err
	}

	if err := urlval.Decode(r.URL.Query(), &dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request URL params")
	}

	return &dst, dst.Validate()
}
