package service

import (
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
			handlers.CtxOrdersQ(s.cfg.DB()),
			handlers.CtxBlockQ(s.cfg.DB()),
			handlers.CtxMatchOrdersQ(s.cfg.DB()),
		),
	)
	r.Route("/integrations/order-aggregator", func(r chi.Router) {
		r.Route("/match_orders", func(r chi.Router) {
			r.Post("/", handlers.AddMatch)
			r.Patch("/", handlers.UpdateMatch)
		})
		r.Route("/orders", func(r chi.Router) {
			r.Post("/", handlers.AddOrder)
			r.Patch("/", handlers.UpdateOrder)
		})
		r.Route("/block", func(r chi.Router) {
			r.Post("/", handlers.SetBlock)
			r.Get("/", handlers.GetBlock)
		})
	})

	return r
}
