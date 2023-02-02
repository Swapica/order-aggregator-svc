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
	a := r.Data.Attributes
	originOrder, originChain := &r.Data.Relationships.OriginOrder, &r.Data.Relationships.OriginChain
	return val.Errors{
		"data/id":                                  val.Validate(r.Data.ID, val.Empty),
		"data/type":                                val.Validate(r.Data.Type, val.Required, val.In(resources.MATCH_ORDER)),
		"data/attributes/match_id":                 val.Validate(a.MatchId, val.Required, val.Min(0)),
		"data/attributes/src_chain":                val.Validate(a.SrcChain, val.Required, val.Min(1)),
		"data/attributes/account":                  val.Validate(a.Account, val.Required, val.Match(addressRegexp)),
		"data/attributes/tokenToSell":              val.Validate(a.TokenToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/amountToSell":             validateUint(a.AmountToSell, amountBitSize),
		"data/attributes/state":                    val.Validate(a.State, val.Required),
		"data/relationships/originOrder/data/id":   validateUint(safeGetKey("id", originOrder), bigintBitSize),
		"data/relationships/originOrder/data/type": val.Validate(safeGetKey("type", originOrder), val.Required, val.In(resources.ORDER)),
		"data/relationships/originChain/data/id":   validateUint(safeGetKey("id", originChain), bigintBitSize),
		"data/relationships/originChain/data/type": val.Validate(safeGetKey("type", originChain), val.Required, val.In(resources.CHAIN)),
	}.Filter()
}

func (r *AddMatch) DBModel() data.Match {
	return data.Match{
		ID:           mustParseBigint(r.Data.ID),
		SrcChain:     *r.Data.Attributes.SrcChain,
		MatchID:      *r.Data.Attributes.MatchId,
		OrderID:      mustParseBigint(r.Data.Relationships.OriginOrder.Data.ID),
		OrderChain:   mustParseBigint(r.Data.Relationships.OriginChain.Data.ID),
		Account:      r.Data.Attributes.Account,
		TokenToSell:  r.Data.Attributes.TokenToSell,
		AmountToSell: r.Data.Attributes.AmountToSell,
		State:        r.Data.Attributes.State,
	}
}
