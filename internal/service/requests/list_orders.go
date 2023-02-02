package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
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
		return nil, toDecodeErr(err, "query")
	}

	return &dst, dst.Validate()
}
