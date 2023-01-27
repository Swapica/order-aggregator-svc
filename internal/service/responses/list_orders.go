package responses

import (
	"math/big"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/pkg/errors"
)

func NewOrderListResponse(orders []data.Order) resources.OrderListResponse {
	list := make([]resources.Order, len(orders))
	for i, o := range orders {
		list[i] = newOrderResource(o)
	}
	return resources.OrderListResponse{Data: list}
}

func newOrderResource(o data.Order) resources.Order {
	var matchSw *string
	if o.MatchSwapica.String != "" {
		matchSw = &o.MatchSwapica.String
	}

	return resources.Order{
		Key: resources.Key{
			ID:   o.ID,
			Type: resources.ORDER,
		},
		Attributes: resources.OrderAttributes{
			Account:      o.Account,
			AmountToBuy:  parseBig(o.AmountToBuy, "amountToBuy"),
			AmountToSell: parseBig(o.AmountToSell, "amountToSell"),
			DestChain:    parseBig(o.DestChain, "destChain"),
			ExecutedBy:   parseBig(o.ExecutedBy.String, "executedBy"),
			MatchSwapica: matchSw,
			State:        o.State,
			TokenToBuy:   o.TokenToBuy,
			TokenToSell:  o.TokenToSell,
		},
	}
}

func parseBig(value, field string) *big.Int {
	if value == "" {
		return nil
	}
	res, ok := new(big.Int).SetString(value, 10)
	if !ok {
		panic(errors.Errorf("failed to parse big.Int from DB string field: %s=%s", field, value))
	}
	return res
}
