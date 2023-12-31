package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/internal/utils"
	"github.com/0loff/grade_gophermart/user"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	usecase user.UseCase
}

func NewAuthMiddleware(usecase user.UseCase) *AuthMiddleware {
	return (&AuthMiddleware{
		usecase: usecase,
	})
}

func (m *AuthMiddleware) Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthCookie, err := r.Cookie("Auth")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				logger.Log.Error("Authentication cookies were not set", zap.Error(err))
				w.WriteHeader(http.StatusUnauthorized)
				return

			default:
				logger.Log.Error("Internal server error. Can't get auth cookie from request.", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		uid, err := m.usecase.ParseToken(r.Context(), AuthCookie.Value)
		if err != nil {
			logger.Log.Error("Failed to get user id from token", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), utils.ContextKeyUID, uid))

		h.ServeHTTP(w, r)
	})
}
