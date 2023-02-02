package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
)

type ListOrders struct {
	page.CursorParams
	FilterChain     *int64  `filter:"chain"`
	FilterBuyToken  *string `filter:"tokenToBuy"`
	FilterSellToken *string `filter:"tokenToSell"`
	FilterAccount   *string `filter:"account"`
	FilterState     *uint8  `filter:"state"`
}

func NewListOrders(r *http.Request) (*ListOrders, error) {
	var dst ListOrders
	if err := urlval.Decode(r.URL.Query(), &dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request URL params")
	}

	return &dst, dst.Validate()
}
