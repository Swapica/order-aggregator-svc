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
	a := r.Data.Attributes
	executedBy, destChain := r.Data.Relationships.ExecutedBy, &r.Data.Relationships.DestChain
	return val.Errors{
		"data/id":                                 val.Validate(r.Data.ID, val.Empty),
		"data/type":                               val.Validate(r.Data.Type, val.Required, val.In(resources.ORDER)),
		"data/attributes/order_id":                val.Validate(a.OrderId, val.Required, val.Min(0)),
		"data/attributes/src_chain":               val.Validate(a.SrcChain, val.Required, val.Min(1)),
		"data/attributes/account":                 val.Validate(a.Account, val.Required, val.Match(addressRegexp)),
		"data/attributes/tokenToSell":             val.Validate(a.TokenToSell, val.Required, val.Match(addressRegexp)),
		"data/attributes/tokenToBuy":              val.Validate(a.TokenToBuy, val.Required, val.Match(addressRegexp)),
		"data/attributes/amountToSell":            validateUint(a.AmountToSell, amountBitSize),
		"data/attributes/amountToBuy":             validateUint(a.AmountToBuy, amountBitSize),
		"data/attributes/state":                   val.Validate(a.State, val.Required, val.Min(uint8(1))),
		"data/attributes/matchSwapica":            val.Validate(a.MatchSwapica, val.NilOrNotEmpty, val.Match(addressRegexp)),
		"data/relationships/executedBy/data/id":   validateOptionalUint(safeGetKey(executedBy).ID, bigintBitSize),
		"data/relationships/executedBy/data/type": val.Validate(safeGetKey(executedBy).Type, val.In(resources.MATCH_ORDER)),
		"data/relationships/destChain/data/id":    validateUint(safeGetKey(destChain).ID, bigintBitSize),
		"data/relationships/destChain/data/type":  val.Validate(safeGetKey(destChain).Type, val.Required, val.In(resources.CHAIN)),
	}.Filter()
}

func (r *AddOrder) DBModel() data.Order {
	execBy := safeGetKey(r.Data.Relationships.ExecutedBy).ID
	matchSw := ""
	if ptr := r.Data.Attributes.MatchSwapica; ptr != nil {
		matchSw = *ptr
	}

	return data.Order{
		SrcChain:     *r.Data.Attributes.SrcChain,
		OrderID:      *r.Data.Attributes.OrderId,
		Account:      r.Data.Attributes.Account,
		TokenToSell:  r.Data.Attributes.TokenToSell,
		TokenToBuy:   r.Data.Attributes.TokenToBuy,
		AmountToSell: r.Data.Attributes.AmountToSell,
		AmountToBuy:  r.Data.Attributes.AmountToBuy,
		DestChain:    mustParseBigint(r.Data.Relationships.DestChain.Data.ID),
		State:        r.Data.Attributes.State,
		ExecutedBy:   sql.NullString{String: execBy, Valid: execBy != ""},
		MatchSwapica: sql.NullString{String: execBy, Valid: matchSw != ""},
	}
}
