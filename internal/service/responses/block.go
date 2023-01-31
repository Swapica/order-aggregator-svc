package responses

import (
	"strconv"

	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewBlockResponse(id uint64) resources.BlockResponse {
	return resources.BlockResponse{
		Data: resources.Block{
			Key: resources.Key{
				ID:   strconv.FormatUint(id, 10),
				Type: resources.BLOCK,
			},
		},
	}
}
