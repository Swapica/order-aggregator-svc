package responses

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func ToTokenResource(t data.Token) resources.Token {
	return resources.Token{
		Key: resources.NewKeyInt64(t.ID, resources.TOKEN),
		Attributes: resources.TokenAttributes{
			Address:  t.Address,
			Decimals: t.Decimals,
			Name:     t.Name,
			SrcChain: t.SrcChain,
			Symbol:   t.Symbol,
		},
	}
}
