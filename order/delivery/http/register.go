package http

import (
	"github.com/0loff/grade_gophermart/order"
	"github.com/go-chi/chi/v5"
)

func RegisterHTTPEndpoints(r chi.Router, uc order.UseCase) {
	h := NewHandler(uc)

	r.Route("/api/user/orders", func(r chi.Router) {
		r.Post("/", h.SetOrder)
		r.Get("/", h.GetOrders)
	})
}
