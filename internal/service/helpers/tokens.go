package helpers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/internal/service/helpers/erc20"
	"github.com/Swapica/order-aggregator-svc/resources"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
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

	endpoint := chain.Attributes.ChainParams.Rpc
	if endpoint == nil {
		return TokenMetadata{}, errors.Errorf(
			"unable to call network for chain=%s: data/attributes/chain_params/rpc is nil", chain.ID)
	}

	cli, err := ethclient.Dial(*endpoint)
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

func IsBadTokenErr(err error) bool {
	// TODO: test it for the different situations
	// 1) address is not a contract
	// 2) address has no name(), symbol() or decimals() method
	// 3) RPC provider errors, like 'invalid project id' (this may be detected as 400, but must be 500;
	// to fix it, you can try calling this method inside fetchTokenMetadata)
	err = errors.Cause(err)
	if e, ok := err.(rpc.HTTPError); ok {
		return e.StatusCode >= http.StatusBadRequest &&
			e.StatusCode < http.StatusInternalServerError
	}

	return err == bind.ErrNoCode
}
