package requests

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type SetBlock struct {
	Number uint64
	Chain  string
}

func NewSetBlock(r *http.Request) (*SetBlock, error) {
	var dst resources.BlockResponse
	if err := json.NewDecoder(r.Body).Decode(&dst); err != nil {
		return nil, errors.Wrap(err, "failed to decode request body")
	}
	num, err := strconv.ParseUint(dst.Data.ID, 10, 64)
	if err != nil {
		return nil, val.Errors{"data/id": errors.Wrap(err, "failed to parse block number")}
	}

	req := SetBlock{
		Number: num,
		Chain:  chi.URLParam(r, "chain"),
	}
	return &req, val.Errors{
		"{chain}":   validateUint(req.Chain, bigintBitSize),
		"data/type": val.Validate(dst.Data.Type, val.Required, val.In(resources.BLOCK)),
	}.Filter()
}
