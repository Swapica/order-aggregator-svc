package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
)

type GetBlockRequest struct {
	Chain string
}

func NewGetBlockRequest(r *http.Request) (GetBlockRequest, error) {
	dst := GetBlockRequest{Chain: chi.URLParam(r, "chain")}
	return dst, val.Errors{
		"{chain}": val.Validate(dst.Chain, val.Required),
	}.Filter()
}
