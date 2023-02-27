package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
	"gitlab.com/distributed_lab/urlval"
)

type ListMatches struct {
	page.Params
	FilterSrcChain         *int64  `filter:"src_chain"`
	FilterState            *uint8  `filter:"state"`
	FilterCreator          *string `filter:"creator"`
	FilterExpired          *bool   `filter:"expired"`
	IncludeSrcChain        bool    `include:"src_chain"`
	IncludeOriginChain     bool    `include:"origin_chain"`
	IncludeOriginOrder     bool    `include:"origin_order"`
	IncludeOriginBuyToken  bool    `include:"origin_order.token_to_buy"`
	IncludeOriginSellToken bool    `include:"origin_order.token_to_sell"`
	IncludeSellToken       bool    `include:"token_to_sell"`
}

func NewListMatches(r *http.Request) (*ListMatches, error) {
	var dst ListMatches
	if err := urlval.Decode(r.URL.Query(), &dst); err != nil {
		return nil, toDecodeErr(err, "query")
	}

	return &dst, dst.Validate()
}
