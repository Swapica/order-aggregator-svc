package helpers

import (
	"context"
	"strings"
	"time"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/internal/service/helpers/erc20"
	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TokenMetadata struct {
	Name, Symbol string
	Decimals     uint8
}

const nativeTokenAddress = "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
const nativeTokenName = "Native Currency"

func GetTokenMetadata(address string, chain resources.Chain) (TokenMetadata, error) {
	if strings.Contains(strings.ToLower(address), nativeTokenAddress) {
		return TokenMetadata{
			Name:     nativeTokenName,
			Symbol:   chain.Attributes.ChainParams.NativeSymbol,
			Decimals: chain.Attributes.ChainParams.NativeDecimals,
		}, nil
	}

	rpc := chain.Attributes.ChainParams.Rpc
	if rpc == nil {
		return TokenMetadata{}, errors.Errorf(
			"unable to call network for chain=%s: data/attributes/chain_params/rpc is nil", chain.ID)
	}

	cli, err := ethclient.Dial(*rpc)
	if err != nil {
		return TokenMetadata{}, errors.Wrap(err, "failed to connect EVM network")
	}
	caller, err := erc20.NewERC20Caller(common.HexToAddress(address), cli)
	if err != nil {
		return TokenMetadata{}, errors.Wrap(err, "failed to create ERC20 contract caller")
	}

	return fetchTokenMetadata(caller)
}

func (m TokenMetadata) DBModel(address string, srcChain int64) data.Token {
	return data.Token{
		Address:  address,
		SrcChain: srcChain,
		Name:     m.Name,
		Symbol:   m.Symbol,
		Decimals: m.Decimals,
	}
}

func fetchTokenMetadata(caller *erc20.ERC20Caller) (TokenMetadata, error) {
	var res TokenMetadata
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := &bind.CallOpts{Context: ctx}

	if res.Name, err = caller.Name(opts); err != nil {
		return res, errors.Wrap(err, "failed to fetch ERC20 name")
	}

	if res.Symbol, err = caller.Symbol(opts); err != nil {
		return res, errors.Wrap(err, "failed to fetch ERC20 symbol")
	}

	res.Decimals, err = caller.Decimals(opts)
	return res, errors.Wrap(err, "failed to fetch ERC20 decimals")
}
