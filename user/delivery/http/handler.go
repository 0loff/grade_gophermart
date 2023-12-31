package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/user"
	"go.uber.org/zap"
)

type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type Handler struct {
	useCase user.UseCase
}

func NewHandler(useCase user.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	credentials := new(Credentials)
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("Cannot decode request JSON body", zap.Error(err))
		return
	}

	token, err := h.useCase.SignUp(r.Context(), credentials.Username, credentials.Password)
	if err != nil {
		if errors.Is(err, user.ErrUserName) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Auth",
		Value: token,
		Path:  "/",
	})

	w.WriteHeader(http.StatusOK)
}

type signInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	credentials := new(Credentials)
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("Cannot decode request JSON body", zap.Error(err))
		return
	}

	token, err := h.useCase.SignIn(r.Context(), credentials.Username, credentials.Password)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrInvalidPassword) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Auth",
		Value: token,
		Path:  "/",
	})

	w.WriteHeader(http.StatusOK)
}
