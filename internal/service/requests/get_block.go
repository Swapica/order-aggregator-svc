package requests

import (
	"net/http"
	"strings"

	val "github.com/go-ozzo/ozzo-validation/v4"
)

type GetBlockRequest struct {
	Chain string
}

func NewGetBlockRequest(r *http.Request) (GetBlockRequest, error) {
	dst := GetBlockRequest{
		Chain: strings.Join(r.URL.Query()["filter[chain]"], ""),
	}
	return dst, val.Errors{
		"filter[chain]": val.Validate(dst.Chain, val.Required),
	}.Filter()
}
