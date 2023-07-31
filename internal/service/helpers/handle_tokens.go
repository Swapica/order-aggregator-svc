package helpers

import (
	"math"
	"math/big"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func GetOrAddToken(q data.Tokens, address string, srcChain resources.Chain) (data.Token, error) {
	chainID := srcChain.Attributes.ChainParams.ChainId
	token, err := q.FilterByAddress(address).FilterByChain(chainID).Get()
	if err != nil {
		return data.Token{}, errors.Wrap(err, "failed to get token by address")
	}

	if token == nil {
		md, err := GetTokenMetadata(address, srcChain)
		if err != nil {
			return data.Token{}, errors.Wrap(err, "failed to get metadata of the token to sell")
		}

		dbt, err := q.Insert(md.DBModel(address, chainID))
		return dbt, errors.Wrap(err, "failed to add token")
	}
	return *token, nil
}

func ConvertAmount(wei *big.Int, decimals uint8) *big.Float {
	ether := new(big.Float)
	weiFloat := new(big.Float).SetInt(wei)
	decimal := new(big.Float).SetInt(big.NewInt(int64(math.Pow10(int(decimals)))))
	ether.Quo(weiFloat, decimal)

	return ether
}
