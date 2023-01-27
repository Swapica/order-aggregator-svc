package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type ListRequest struct {
	Chain string // fixme: skip all tags of urlval?
	page.CursorParams
}

func NewListRequest(r *http.Request) (*ListRequest, error) {
	dst := ListRequest{Chain: chi.URLParam(r, "chain")}
	if err := urlval.Decode(r.URL.Query(), &dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request URL params")
	}

	err := val.Errors{"{chain}": val.Validate(dst.Chain, val.Required)}.Filter()
	if err != nil {
		return nil, err
	}

	return &dst, dst.Validate()
}
