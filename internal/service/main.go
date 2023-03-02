package service

import (
	"net/http"

	"github.com/Swapica/order-aggregator-svc/internal/config"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log *logan.Entry
	cfg config.Config
}

func (s *service) run() error {
	s.log.Info("Service started")
	r := s.router()

	if err := s.cfg.Copus().RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.cfg.Listener(), r)
}

func newService(cfg config.Config) *service {
	return &service{
		log: cfg.Log(),
		cfg: cfg,
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(); err != nil {
		panic(err)
	}
}
