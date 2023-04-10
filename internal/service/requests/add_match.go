package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

type AddMatch resources.AddMatchRequest

func NewAddMatch(r *http.Request) (*AddMatch, error) {
	var dst AddMatch
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, toDecodeErr(err, "body")
	}

	return &dst, dst.validate()
}

func (r *AddMatch) validate() error {
	a := r.Data.Attributes
	return val.Errors{
		"data/id":                         val.Validate(r.Data.ID, val.Empty),
		"data/type":                       val.Validate(r.Data.Type, val.Required, val.In(resources.MATCH_ORDER)),
		"data/attributes/match_id":        val.Validate(a.MatchId, val.Required, val.Min(1)),
		"data/attributes/creator":         val.Validate(a.Creator, val.Required, val.Match(addressRegexp)),
		"data/attributes/token_to_sell":   val.Validate(a.TokenToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/amount_to_sell":  validateUint(a.AmountToSell, amountBitSize),
		"data/attributes/state":           val.Validate(a.State, val.Required, val.In(data.StateAwaitingFinalization)),
		"data/attributes/src_chain_id":    val.Validate(a.SrcChainId, val.Required, val.Min(1)),
		"data/attributes/origin_chain_id": val.Validate(a.OriginChainId, val.Required, val.Min(1)),
		"data/attributes/origin_order_id": val.Validate(a.OriginOrderId, val.Required, val.Min(1)),
		"data/attributes/auto_execute":    val.Validate(a.AutoExecute, val.NotNil),
	}.Filter()
}

func (r *AddMatch) DBModel(originOrder, sellToken int64) data.Match {
	return data.Match{
		MatchID:     r.Data.Attributes.MatchId,
		SrcChain:    r.Data.Attributes.SrcChainId,
		OriginOrder: originOrder,
		OrderID:     r.Data.Attributes.OriginOrderId,
		OrderChain:  r.Data.Attributes.OriginChainId,
		Creator:     r.Data.Attributes.Creator,
		SellToken:   sellToken,
		SellAmount:  r.Data.Attributes.AmountToSell,
		State:       r.Data.Attributes.State,
		AutoExecute: r.Data.Attributes.AutoExecute,
	}
}
