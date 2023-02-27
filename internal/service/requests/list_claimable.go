package requests

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/service/page"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval"
)

type ListClaimable struct {
	page.Params
	FilterSrcChain         *int64  `filter:"src_chain"`
	FilterCreator          *string `filter:"creator"`
	IncludeSrcChain        bool    `include:"src_chain"`
	IncludeOriginChain     bool    `include:"origin_chain"`
	IncludeOriginBuyToken  bool    `include:"origin_order.token_to_buy"`
	IncludeOriginSellToken bool    `include:"origin_order.token_to_sell"`
	IncludeSellToken       bool    `include:"token_to_sell"`
}

func NewListClaimable(r *http.Request) (*ListClaimable, error) {
	var dst ListClaimable
	if err := urlval.Decode(r.URL.Query(), &dst); err != nil {
		return nil, toDecodeErr(err, "query")
	}
	if dst.FilterCreator == nil {
		return nil, val.Errors{"filter[creator]": val.ErrRequired}
	}

	return &dst, dst.Validate()
}
