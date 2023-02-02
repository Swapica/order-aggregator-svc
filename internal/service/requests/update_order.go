package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateOrder struct {
	Body           resources.UpdateOrderRequest
	OrderID, Chain int64
	ExecutedBy     *int64
}

func NewUpdateOrder(r *http.Request) (*UpdateOrder, error) {
	var dst UpdateOrder
	if err := json.NewDecoder(r.Body).Decode(&dst.Body); err != nil {
		return nil, toDecodeErr(err, "body")
	}

	var errChain, errOrderID, errExecutedBy error
	dst.Chain, errChain = parseBigint(chi.URLParam(r, "chain"))
	dst.OrderID, errOrderID = parseBigint(dst.Body.Data.ID)

	if rel := dst.Body.Data.Relationships; rel != nil {
		var ex int64
		ex, errExecutedBy = parseBigint(safeGetKey(rel.ExecutedBy).ID)
		dst.ExecutedBy = &ex
	}

	a := dst.Body.Data.Attributes
	return &dst, val.Errors{
		"{chain}":                      errChain,
		"data/id":                      errOrderID,
		"data/type":                    val.Validate(dst.Body.Data.Type, val.Required, val.In(resources.ORDER)),
		"data/attributes/state":        val.Validate(a.State, val.Required, val.Min(uint8(1))),
		"data/attributes/executedBy":   errExecutedBy,
		"data/attributes/matchSwapica": val.Validate(a.MatchSwapica, val.NilOrNotEmpty, val.Match(addressRegexp)),
	}.Filter()
}
