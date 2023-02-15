package mem

import (
	"github.com/Swapica/order-aggregator-svc/internal/data"
	"github.com/Swapica/order-aggregator-svc/resources"
)

func NewChains(chains []resources.Chain) data.Chains {
	return &chainsQ{
		chains:  chains,
		filters: make([]chainFilter, 0),
	}
}

type chainsQ struct {
	chains  []resources.Chain
	filters []chainFilter
}

type chainFilter func(value resources.Chain) bool

func (q *chainsQ) New() data.Chains {
	return NewChains(q.chains)
}

func (q *chainsQ) Get() *resources.Chain {
	for _, value := range q.chains {
		if q.filter(value) {
			return &value
		}
	}

	return nil
}

func (q *chainsQ) Select() []resources.Chain {
	result := make([]resources.Chain, 0, len(q.chains))
	for _, value := range q.chains {
		// duplicates are OK, because resources.Included.Add will eliminate them
		if q.filter(value) {
			result = append(result, value)
		}
	}

	return result
}

func (q *chainsQ) FilterByID(ids ...string) data.Chains {
	q.filters = append(q.filters, func(value resources.Chain) bool {
		return contains(ids, value.ID)
	})
	return q
}

func (q *chainsQ) FilterByChainID(ids ...int64) data.Chains {
	q.filters = append(q.filters, func(value resources.Chain) bool {
		return contains(ids, value.Attributes.ChainParams.ChainId)
	})
	return q
}

func (q *chainsQ) filter(value resources.Chain) bool {
	for _, filter := range q.filters {
		if !filter(value) {
			return false
		}
	}

	return true
}

func contains[T comparable](src []T, value T) bool {
	for _, v := range src {
		if v == value {
			return true
		}
	}

	return false
}
