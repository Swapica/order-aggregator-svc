package service

import (
	"github.com/Swapica/order-aggregator-svc/internal/data/postgres"
	"github.com/Swapica/order-aggregator-svc/internal/service/handlers"
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
)

func (s *service) router() chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			handlers.CtxLog(s.log),
			handlers.CtxOrdersQ(postgres.NewOrders(s.cfg.DB())),
			handlers.CtxMatchOrdersQ(postgres.NewMatchOrders(s.cfg.DB())),
			handlers.CtxBlockQ(postgres.NewLastBlock(s.cfg.DB())),
		),
	)
	r.Route("/integrations/order-aggregator", func(r chi.Router) {
		r.Route("/{chain}", func(r chi.Router) {
			r.Route("/match_orders", func(r chi.Router) {
				r.Post("/", handlers.AddMatch)
				r.Patch("/", handlers.UpdateMatch)
				r.Get("/", handlers.ListMatches)
			})
			r.Route("/orders", func(r chi.Router) {
				r.Post("/", handlers.AddOrder)
				r.Patch("/", handlers.UpdateOrder)
				r.Get("/", handlers.ListOrders)
			})
			r.Route("/block", func(r chi.Router) {
				r.Post("/", handlers.SetBlock)
				r.Get("/", handlers.GetBlock)
			})
		})
	})

	return r
}
