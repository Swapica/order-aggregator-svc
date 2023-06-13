package config

import (
	"crypto/ecdsa"
	figure "gitlab.com/distributed_lab/figure/v3"

	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Notifications struct {
	PushURL         string            `fig:"push_url,required"`
	ChannelAddress  string            `fig:"channel_address,required"`
	PushCommAddress string            `fig:"push_comm_address,required"`
	PrivateKey      *ecdsa.PrivateKey `fig:"private_key,required"`
}

func (c *config) Notifications() Notifications {
	return c.notifications.Do(func() interface{} {
		var result Notifications

		err := figure.
			Out(&result).
			With(figure.BaseHooks, figure.EthereumHooks).
			From(kv.MustGetStringMap(c.getter, "notifications")).
			Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out ethereum config"))
		}

		return result
	}).(Notifications)
}
