package requests

import (
	"net/http"

	"github.com/go-chi/chi"
)

type GetBlockRequest struct {
	Chain string
}

func NewGetBlockRequest(r *http.Request) (GetBlockRequest, error) {
	dst := GetBlockRequest{Chain: chi.URLParam(r, "chain")}
	return dst, validateChain(dst.Chain)
}
