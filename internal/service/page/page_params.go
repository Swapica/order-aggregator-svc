package page

import (
	"math"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	pageParamLimit  = "page[limit]"
	pageParamNumber = "page[number]"
	pageParamOrder  = "page[order]"

	maxLimit uint64 = 100
)

type Params struct {
	pgdb.OffsetPageParams
}

func (p *Params) Validate() error {
	return val.Errors{
		pageParamLimit:  val.Validate(p.Limit, val.Max(maxLimit)),
		pageParamOrder:  val.Validate(p.Order, val.In(pgdb.OrderTypeAsc, pgdb.OrderTypeDesc)),
		pageParamNumber: val.Validate(p.PageNumber, val.Max(uint64(math.MaxInt64))),
	}.Filter()
}
