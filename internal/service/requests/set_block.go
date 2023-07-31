package requests

import (
	"encoding/json"
	"net/http"

	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

type SetBlock struct {
	Number, Chain int64
}

func NewSetBlock(r *http.Request) (*SetBlock, error) {
	var dst resources.BlockResponse
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, toDecodeErr(err, "body")
	}

	c, errChain := parseBigint(chi.URLParam(r, "chain"))
	n, errNumber := parseBigint(dst.Data.ID)

	return &SetBlock{Chain: c, Number: n}, val.Errors{
		"{chain}":   errChain,
		"data/id":   errNumber,
		"data/type": val.Validate(dst.Data.Type, val.Required, val.In(resources.BLOCK)),
	}.Filter()
}
