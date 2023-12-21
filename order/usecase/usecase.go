package usecase

import (
	"context"
	"errors"

	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/models"
	"github.com/0loff/grade_gophermart/order"
	"github.com/ShiraazMoollatjie/goluhn"
	"go.uber.org/zap"
)

type OrderUseCase struct {
	orderRepo order.Repository
}

func NewOrderUseCase(
	orderRepo order.Repository,
) *OrderUseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
	}
}

func (o OrderUseCase) SetOrder(ctx context.Context, Order models.Order) error {
	err := goluhn.Validate(Order.OrderNum)
	if err != nil || len(Order.OrderNum) < 3 {
		logger.Log.Error("Order number is incorrect", zap.Error(err))
		return order.ErrInvalidOrderNumber
	}

	if err := o.orderRepo.InsertOrder(ctx, Order); err != nil {
		switch {
		case errors.Is(err, order.ErrOrderBeenRegistered):
			orderUUID, err := o.orderRepo.GetUUIDByOrder(ctx, Order.OrderNum)
			if err != nil {
				logger.Log.Error("")
				return err
			}

			if orderUUID != Order.UUID {
				logger.Log.Error("The order belongs to another user")
				return order.ErrAnotherUserOrder
			}

			return order.ErrCurrentUserOrder

		default:
			logger.Log.Error("Internal server error", zap.Error(err))
			return err
		}
	}

	return nil
}

func (o OrderUseCase) GetUserOrders(ctx context.Context, uuid string) ([]models.Order, error) {
	return o.orderRepo.GetOrdersByUUID(ctx, uuid)
}
