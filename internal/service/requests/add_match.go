package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

type AddMatch resources.MatchResponse

func NewAddMatch(r *http.Request) (*AddMatch, error) {
	var dst AddMatch
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, toDecodeErr(err, "body")
	}

	return &dst, dst.validate()
}

func (r *AddMatch) validate() error {
	a, rel := r.Data.Attributes, r.Data.Relationships
	return val.Errors{
		"data/id":                                   val.Validate(r.Data.ID, val.Empty),
		"data/type":                                 val.Validate(r.Data.Type, val.Required, val.In(resources.MATCH_ORDER)),
		"data/attributes/match_id":                  val.Validate(a.MatchId, val.Required, val.Min(0)),
		"data/attributes/creator":                   val.Validate(a.Creator, val.Required, val.Match(addressRegexp)),
		"data/attributes/token_to_sell":             val.Validate(a.TokenToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/amount_to_sell":            validateUint(a.AmountToSell, amountBitSize),
		"data/attributes/state":                     val.Validate(a.State, val.Required, val.In(data.StateAwaitingFinalization)),
		"data/relationships/src_chain/data/id":      validateUint(safeGetKey(&rel.SrcChain).ID, bigintBitSize),
		"data/relationships/src_chain/data/type":    val.Validate(safeGetKey(&rel.SrcChain).Type, val.Required, val.In(resources.CHAIN)),
		"data/relationships/origin_chain/data/id":   validateUint(safeGetKey(&rel.OriginChain).ID, bigintBitSize),
		"data/relationships/origin_chain/data/type": val.Validate(safeGetKey(&rel.OriginChain).Type, val.Required, val.In(resources.CHAIN)),
		"data/relationships/origin_order/data/id":   validateUint(safeGetKey(&rel.OriginOrder).ID, bigintBitSize),
		"data/relationships/origin_order/data/type": val.Validate(safeGetKey(&rel.OriginOrder).Type, val.Required, val.In(resources.ORDER)),
	}.Filter()
}

func (r *AddMatch) DBModel() data.Match {
	return data.Match{
		SrcChain:   mustParseBigint(r.Data.Relationships.SrcChain.Data.ID),
		MatchID:    *r.Data.Attributes.MatchId,
		OrderID:    mustParseBigint(r.Data.Relationships.OriginOrder.Data.ID),
		OrderChain: mustParseBigint(r.Data.Relationships.OriginChain.Data.ID),
		Creator:    r.Data.Attributes.Creator,
		SellToken:  r.Data.Attributes.TokenToSell,
		SellAmount: r.Data.Attributes.AmountToSell,
		State:      r.Data.Attributes.State,
	}
}
