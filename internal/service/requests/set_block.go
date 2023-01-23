package requests

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type SetBlockRequest struct {
	Number uint64
	Chain  string
}

func NewSetBlockRequest(r *http.Request) (*SetBlockRequest, error) {
	var dst resources.BlockRequest
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}
	num, err := strconv.ParseUint(dst.Data.ID, 10, 64)
	if err != nil {
		return nil, val.Errors{"data/id": errors.Wrap(err, "failed to parse block number")}
	}

	return &SetBlockRequest{
			Number: num,
			Chain:  dst.Data.Attributes.Chain,
		}, val.Errors{
			"data/type":             val.Validate(dst.Data.Type, val.Required, val.In(resources.BLOCK)),
			"data/attributes/chain": val.Validate(dst.Data.Attributes.Chain, val.Required),
		}.Filter()
}
