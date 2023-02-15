package config

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Swapica/order-aggregator-svc/resources"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func (c *config) Chains() []resources.Chain {
	return c.chains.Do(func() interface{} {
		var cfg struct {
			*url.URL `fig:"url,required"`
		}

		err := figure.Out(&cfg).From(kv.MustGetStringMap(c.getter, "chains")).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out chains endpoint"))
		}

		resp, err := http.Get(cfg.URL.String())
		if err != nil {
			panic(errors.Wrap(err, "failed to fetch chain list"))
		}

		var chains resources.ChainListResponse
		if err = json.NewDecoder(resp.Body).Decode(&chains); err != nil {
			panic(errors.Wrap(err, "failed to unmarshal chain list response"))
		}

		return chains.Data
	}).([]resources.Chain)
}
