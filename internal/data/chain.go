package data

import (
	"github.com/Swapica/order-aggregator-svc/resources"
)

type Chains interface {
	New() Chains
	Get() *resources.Chain
	Select() []resources.Chain
	SelectIDs() []int64

	FilterByID(ids ...string) Chains
	FilterByChainID(ids ...int64) Chains
}
