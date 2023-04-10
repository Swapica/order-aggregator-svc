package requests

import (
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
		"data/attributes/order_id":       val.Validate(a.OrderId, val.Required, val.Min(1)),
		"data/attributes/creator":        val.Validate(a.Creator, val.Required, val.Match(addressRegexp)),
		"data/attributes/token_to_sell":  val.Validate(a.TokenToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/token_to_buy":   val.Validate(a.TokenToBuy, val.Required, val.Match(addressRegexp)),
		"data/attributes/amount_to_sell": validateUint(a.AmountToSell, amountBitSize),
		"data/attributes/amount_to_buy":  validateUint(a.AmountToBuy, amountBitSize),
		"data/attributes/state":          val.Validate(a.State, val.Required, val.In(data.StateAwaitingMatch)),
		"data/attributes/src_chain_id":   val.Validate(a.SrcChainId, val.Required, val.Min(1)),
		"data/attributes/dest_chain_id":  val.Validate(a.DestChainId, val.Required, val.Min(1)),
		"data/attributes/match_id":       val.Validate(a.MatchId, val.Nil),
		"data/attributes/match_swapica":  val.Validate(a.MatchSwapica, val.Nil),
		"data/attributes/auto_execute":   val.Validate(a.AutoExecute, val.NotNil),
	}.Filter()
}

func (r *AddOrder) DBModel(sellToken, buyToken int64) data.Order {
	return data.Order{
		OrderID:     r.Data.Attributes.OrderId,
		SrcChain:    r.Data.Attributes.SrcChainId,
		Creator:     r.Data.Attributes.Creator,
		SellToken:   sellToken,
		BuyToken:    buyToken,
		SellAmount:  r.Data.Attributes.AmountToSell,
		BuyAmount:   r.Data.Attributes.AmountToBuy,
		DestChain:   r.Data.Attributes.DestChainId,
		State:       r.Data.Attributes.State,
		AutoExecute: r.Data.Attributes.AutoExecute,
		// ExecutedByMatch, MatchID, MatchSwapica must not appear on the order creation, according to the core contract
	}
}
