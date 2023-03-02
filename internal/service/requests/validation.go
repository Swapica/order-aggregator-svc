package requests

import (
	"math/big"
	"regexp"
	"strconv"

	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	bigintBitSize = 63
	amountBitSize = 256
)

var addressRegexp = regexp.MustCompile("^0x[0-9A-Fa-f]{40}$")
var maxUint256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(amountBitSize), nil)
var bigZero = big.NewInt(0)

func validateUint(value string, bitSize int) error {
	if value == "" {
		return val.ErrRequired
	}

	if bitSize <= 64 {
		_, err := strconv.ParseUint(value, 10, bitSize)
		return err
	}

	parsed, success := new(big.Int).SetString(value, 10)
	if !success {
		return errors.Errorf("failed to parse big integer: string=%s", value)
	}
	if isOutOfRange(parsed, bitSize) {
		return errors.Errorf("parsed value is out of range for uint%d type: parsed=%s", bitSize, parsed.String())
	}

	return nil
}

func isOutOfRange(n *big.Int, bitSize int) bool {
	if bitSize == amountBitSize {
		return n.Cmp(bigZero) == -1 || n.Cmp(maxUint256) == 1
	}

	max := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(bitSize)), nil)
	return n.Cmp(bigZero) == -1 || n.Cmp(max) == 1
}

func toDecodeErr(err error, what string) error {
	return val.Errors{"/": errors.Wrap(err, "failed to decode request "+what)}
}

func parseBigint(value string) (int64, error) {
	n, err := strconv.ParseUint(value, 10, bigintBitSize)
	return int64(n), errors.Wrap(err, "failed to parse 63-bit unsigned integer")
}

// mustParseBigint relies on validateUint: if validation succeeded with bitSize=bigintBitSize for value, no panic will appear
func mustParseBigint(value string) int64 {
	n, err := parseBigint(value)
	if err != nil {
		panic(err)
	}
	return n
}

func safeGetKey(rel *resources.Relation) resources.Key {
	if rel != nil && rel.Data != nil {
		return *rel.Data
	}
	return resources.Key{}
}
