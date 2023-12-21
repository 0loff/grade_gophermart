package http

import (
	"github.com/0loff/grade_gophermart/balance"
	"github.com/go-chi/chi/v5"
)

func RegisterHTTPEndpoints(r chi.Router, uc balance.UseCase) {
	h := NewHandler(uc)

	r.Get("/api/user/withdrawals", h.GetWithdrawals)

	r.Route("/api/user/balance", func(r chi.Router) {
		r.Get("/", h.GetBalance)
		r.Post("/withdraw", h.SetOrderWithdraw)
	})
}
