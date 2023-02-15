package requests

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

type AddOrder resources.OrderResponse

func NewAddOrder(r *http.Request) (*AddOrder, error) {
	var dst AddOrder
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, toDecodeErr(err, "body")
	}

	return &dst, dst.validate()
}

func (r *AddOrder) validate() error {
	a, rel := r.Data.Attributes, r.Data.Relationships
	return val.Errors{
		"data/id":                                        val.Validate(r.Data.ID, val.Empty),
		"data/type":                                      val.Validate(r.Data.Type, val.Required, val.In(resources.ORDER)),
		"data/attributes/order_id":                       val.Validate(a.OrderId, val.Required, val.Min(0)),
		"data/attributes/creator":                        val.Validate(a.Creator, val.Required, val.Match(addressRegexp)),
		"data/attributes/token_to_sell":                  val.Validate(a.TokenToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/token_to_buy":                   val.Validate(a.TokenToBuy, val.Required, val.Match(addressRegexp)),
		"data/attributes/amount_to_sell":                 validateUint(a.AmountToSell, amountBitSize),
		"data/attributes/amount_to_buy":                  validateUint(a.AmountToBuy, amountBitSize),
		"data/attributes/state":                          val.Validate(a.State, val.Required, val.In(data.StateAwaitingMatch)),
		"data/attributes/match_swapica":                  val.Validate(a.MatchSwapica, val.NilOrNotEmpty, val.Match(addressRegexp)),
		"data/relationships/match":                       val.Validate(rel.Match, val.Nil),
		"data/relationships/src_chain/data/id":           validateUint(safeGetKey(&rel.SrcChain).ID, bigintBitSize),
		"data/relationships/src_chain/data/type":         val.Validate(safeGetKey(&rel.SrcChain).Type, val.Required, val.In(resources.CHAIN)),
		"data/relationships/destination_chain/data/id":   validateUint(safeGetKey(&rel.DestinationChain).ID, bigintBitSize),
		"data/relationships/destination_chain/data/type": val.Validate(safeGetKey(&rel.DestinationChain).Type, val.Required, val.In(resources.CHAIN)),
	}.Filter()
}

func (r *AddOrder) DBModel() data.Order {
	matchSw := ""
	if ptr := r.Data.Attributes.MatchSwapica; ptr != nil {
		matchSw = *ptr
	}

	return data.Order{
		SrcChain:     mustParseBigint(r.Data.Relationships.SrcChain.Data.ID),
		OrderID:      *r.Data.Attributes.OrderId,
		Creator:      r.Data.Attributes.Creator,
		SellToken:    r.Data.Attributes.TokenToSell,
		BuyToken:     r.Data.Attributes.TokenToBuy,
		SellAmount:   r.Data.Attributes.AmountToSell,
		BuyAmount:    r.Data.Attributes.AmountToBuy,
		DestChain:    mustParseBigint(r.Data.Relationships.DestinationChain.Data.ID),
		State:        r.Data.Attributes.State,
		MatchID:      sql.NullString{}, // must be empty on creation
		MatchSwapica: sql.NullString{String: matchSw, Valid: matchSw != ""},
	}
}
