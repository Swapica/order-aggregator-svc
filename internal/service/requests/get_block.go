package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

type GetBlock struct {
	Chain int64
}

func NewGetBlock(r *http.Request) (GetBlock, error) {
	c, errChain := parseBigint(chi.URLParam(r, "chain"))

	return GetBlock{Chain: c}, val.Errors{
		"{chain}": errChain,
	}.Filter()
}
