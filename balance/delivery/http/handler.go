package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/0loff/grade_gophermart/balance"
	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/internal/utils"
	"github.com/0loff/grade_gophermart/models"
	"github.com/0loff/grade_gophermart/order"
	"go.uber.org/zap"
)

type OrderWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type Handler struct {
	useCase balance.UseCase
}

func NewHandler(useCase balance.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	uuid, ok := utils.GetCtxUID(r.Context())
	if !ok {
		logger.Log.Error("Cannot get UID from context")
	}

	balance, err := h.useCase.GetUserBalance(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(balance); err != nil {
		logger.Log.Error("Error encoding response data", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SetOrderWithdraw(w http.ResponseWriter, r *http.Request) {
	uuid, ok := utils.GetCtxUID(r.Context())
	if !ok {
		logger.Log.Error("Cannot get UID from context")
	}
	orderWithdraw := new(OrderWithdraw)
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(orderWithdraw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Log.Error("Cannot decode request JSON body", zap.Error(err))
		return
	}

	Order := models.Order{
		OrderNum: orderWithdraw.Order,
		Sum:      orderWithdraw.Sum,
		UUID:     uuid,
	}

	if err := h.useCase.SetOrderWithdraw(r.Context(), Order); err != nil {
		switch {
		case errors.Is(err, order.ErrInvalidOrderNumber):
			w.WriteHeader(http.StatusUnprocessableEntity)
			return

		case errors.Is(err, balance.ErrNotEnoughPoints):
			w.WriteHeader(http.StatusPaymentRequired)
			return

		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	uuid, ok := utils.GetCtxUID(r.Context())
	if !ok {
		logger.Log.Error("Cannot get UID from context")
	}

	drawals, err := h.useCase.GetUserWithdrawals(r.Context(), uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	if len(drawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(drawals); err != nil {
		logger.Log.Error("Error encoding response data", zap.Error(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
