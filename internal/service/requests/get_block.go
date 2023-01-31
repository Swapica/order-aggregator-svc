package requests

import (
	"net/http"

	"github.com/go-chi/chi"
)

type GetBlock struct {
	Chain string
}

func NewGetBlock(r *http.Request) (GetBlock, error) {
	dst := GetBlock{Chain: chi.URLParam(r, "chain")}
	return dst, validateChain(dst.Chain)
}
