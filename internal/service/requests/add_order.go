package requests

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

type AddOrder resources.AddOrderRequest

func NewAddOrder(r *http.Request) (*AddOrder, error) {
	var dst AddOrder
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, toDecodeErr(err, "body")
	}

	return &dst, dst.validate()
}

func (r *AddOrder) validate() error {
	a := r.Data.Attributes
	return val.Errors{
		"data/id":                        val.Validate(r.Data.ID, val.Empty),
		"data/type":                      val.Validate(r.Data.Type, val.Required, val.In(resources.ORDER)),
		"data/attributes/order_id":       val.Validate(a.OrderId, val.Required, val.Min(0)),
		"data/attributes/creator":        val.Validate(a.Creator, val.Required, val.Match(addressRegexp)),
		"data/attributes/token_to_sell":  val.Validate(a.TokenToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/token_to_buy":   val.Validate(a.TokenToBuy, val.Required, val.Match(addressRegexp)),
		"data/attributes/amount_to_sell": validateUint(a.AmountToSell, amountBitSize),
		"data/attributes/amount_to_buy":  validateUint(a.AmountToBuy, amountBitSize),
		"data/attributes/state":          val.Validate(a.State, val.Required, val.In(data.StateAwaitingMatch)),
		"data/attributes/match_swapica":  val.Validate(a.MatchSwapica, val.NilOrNotEmpty, val.Match(addressRegexp)),
		"data/attributes/src_chain_id":   val.Validate(a.SrcChainId, val.Required, val.Min(1)),
		"data/attributes/dest_chain_id":  val.Validate(a.DestChainId, val.Required, val.Min(1)),
	}.Filter()
}

func (r *AddOrder) DBModel() data.Order {
	matchSw := ""
	if ptr := r.Data.Attributes.MatchSwapica; ptr != nil {
		matchSw = *ptr
	}

	return data.Order{
		SrcChain:     *r.Data.Attributes.SrcChainId,
		OrderID:      *r.Data.Attributes.OrderId,
		Creator:      r.Data.Attributes.Creator,
		SellToken:    r.Data.Attributes.TokenToSell,
		BuyToken:     r.Data.Attributes.TokenToBuy,
		SellAmount:   r.Data.Attributes.AmountToSell,
		BuyAmount:    r.Data.Attributes.AmountToBuy,
		DestChain:    *r.Data.Attributes.DestChainId,
		State:        r.Data.Attributes.State,
		MatchID:      sql.NullString{},
		MatchSwapica: sql.NullString{String: matchSw, Valid: matchSw != ""},
	}
}
