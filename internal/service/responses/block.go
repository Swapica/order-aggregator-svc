package responses

import (
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewBlock(id int64) resources.BlockResponse {
	return resources.BlockResponse{
		Data: resources.Block{
			Key: resources.NewKeyInt64(id, resources.BLOCK),
		},
	}
}
