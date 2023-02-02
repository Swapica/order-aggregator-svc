package requests

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/Swapica/order-aggregator-svc/resources"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const (
	bigintBitSize = 63
	amountBitSize = 256
)

var addressRegexp = regexp.MustCompile("^0x[0-9A-Fa-f]{40}$")

func validateUint(value string, bitSize int) error {
	if value == "" {
		return val.ErrRequired
	}
	if strings.ContainsAny(value, "eE") {
		return val.ErrMatchInvalid
	}

	if bitSize <= 64 {
		_, err := strconv.ParseUint(value, 10, bitSize)
		return err
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	if parsed < 0 || math.Mod(parsed, 1.0) != 0 {
		return errors.Errorf("parsed value is not uint%d: parsed=%v", bitSize, parsed)
	}
	if max := math.Pow(2, float64(bitSize)); parsed > max {
		return errors.Errorf("parsed value exceeds maximum: parsed=%v max=%v", parsed, max)
	}

	return nil
}

func validateChain(ch string) error {
	return val.Errors{"{chain}": validateUint(ch, bigintBitSize)}.Filter() // EIP 2294
}

func parseBigint(value string) (int64, error) {
	n, err := strconv.ParseUint(value, 10, bigintBitSize)
	return int64(n), errors.Wrap(err, "failed to parse 63-bit unsigned integer")
}

// MustParseBigint relies on ValidateUint: if validation succeeded with bitSize=bigintBitSize for value, no panic will appear
func mustParseBigint(value string) int64 {
	n, err := parseBigint(value)
	if err != nil {
		panic(err)
	}
	return n
}

func safeGetKey(field string, rel *resources.Relation) string {
	if rel != nil && rel.Data != nil {
		switch field {
		case "id":
			return rel.Data.ID
		case "type":
			return string(rel.Data.Type)
		}
	}
	return ""
}
