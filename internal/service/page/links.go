package page

import (
	"net/http"
	"strconv"

	"github.com/Swapica/order-aggregator-svc/resources"
)

func (p *Params) GetLinks(r *http.Request) *resources.Links {
	result := resources.Links{
		Next: p.getOffsetLink(r, p.PageNumber+1),
		Self: p.getOffsetLink(r, p.PageNumber),
	}
	if p.PageNumber > 0 {
		result.Prev = p.getOffsetLink(r, p.PageNumber-1)
	}
	return &result
}

func (p *Params) getOffsetLink(r *http.Request, number uint64) string {
	u := r.URL
	query := u.Query()
	query.Set(pageParamNumber, strconv.FormatUint(number, 10))
	query.Set(pageParamLimit, strconv.FormatUint(p.Limit, 10))
	query.Set(pageParamOrder, p.Order)
	u.RawQuery = query.Encode()
	return u.String()
}
