package http

import (
	"github.com/0loff/grade_gophermart/user"
	"github.com/go-chi/chi/v5"
)

func RegisterHTTPEndpoints(r chi.Router, uc user.UseCase) {
	h := NewHandler(uc)

	r.Post("/api/user/register", h.SignUp)
	r.Post("/api/user/login", h.SignIn)
}
