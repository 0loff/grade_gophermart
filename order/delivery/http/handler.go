package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/internal/utils"
	"github.com/0loff/grade_gophermart/models"
	"github.com/0loff/grade_gophermart/order"
	"go.uber.org/zap"
)

type Handler struct {
	useCase order.UseCase
}

func NewHandler(useCase order.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) SetOrder(w http.ResponseWriter, r *http.Request) {
	uuid, ok := utils.GetCtxUID(r.Context())
	if !ok {
		logger.Log.Error("Cannot get UID from context")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("Error parsing request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Order := &models.Order{
		OrderNum: string(body),
		UUID:     uuid,
	}

	if err = h.useCase.SetOrder(r.Context(), *Order); err != nil {
		switch {
		case errors.Is(err, order.ErrInvalidOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
			return

		case errors.Is(err, order.ErrAnotherUserOrder):
			w.WriteHeader(http.StatusConflict)
			return

		case errors.Is(err, order.ErrCurrentUserOrder):
			w.WriteHeader(http.StatusOK)
			return

		default:
			logger.Log.Error("Internal server error. Unable to register order", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	uuid, ok := utils.GetCtxUID(r.Context())
	if !ok {
		logger.Log.Error("Cannot get UID from context")
	}

	orders, err := h.useCase.GetUserOrders(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(orders); err != nil {
		logger.Log.Error("Error encoding response data", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
