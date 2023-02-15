package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval"
)

type ListOrders struct {
	page.CursorParams
	FilterChain      *int64  `filter:"chain"`
	FilterCreator    *string `filter:"creator"`
	FilterBuyToken   *string `filter:"token_to_buy"`
	FilterSellToken  *string `filter:"token_to_sell"`
	FilterDestChain  *int64  `filter:"destination_chain"`
	FilterState      *uint8  `filter:"state"`
	IncludeSrcChain  bool    `include:"src_chain"`
	IncludeDestChain bool    `include:"destination_chain"`
}

func NewListOrders(r *http.Request) (*ListOrders, error) {
	var dst ListOrders
	if err := urlval.Decode(r.URL.Query(), &dst); err != nil {
		return nil, toDecodeErr(err, "query")
	}

	return &dst, dst.Validate()
}
