package usecase

import (
	"context"

	"github.com/0loff/grade_gophermart/balance"
	"github.com/0loff/grade_gophermart/internal/logger"
	"github.com/0loff/grade_gophermart/models"
	"github.com/0loff/grade_gophermart/order"
	"github.com/ShiraazMoollatjie/goluhn"
	"go.uber.org/zap"
)

type BalanceUseCase struct {
	orderRepo order.Repository
}

func NewBalanceUseCase(
	orderRepo order.Repository,
) *BalanceUseCase {
	return &BalanceUseCase{
		orderRepo: orderRepo,
	}
}

func (b BalanceUseCase) GetUserBalance(ctx context.Context, uuid string) (models.Balance, error) {
	Balance, err := b.orderRepo.GetBalance(ctx, uuid)
	if err != nil {
		logger.Log.Error("impossible to get user balance", zap.Error(err))
	}

	return Balance, err
}

func (b BalanceUseCase) SetOrderWithdraw(ctx context.Context, Order models.Order) error {
	err := goluhn.Validate(Order.OrderNum)
	if err != nil || len(Order.OrderNum) < 3 {
		logger.Log.Error("Order number is incorrect", zap.Error(err))
		return order.ErrInvalidOrderNumber
	}

	Balance, err := b.GetUserBalance(ctx, Order.UUID)
	if err != nil {
		logger.Log.Error("Cannot calculate current user balance", zap.Error(err))
		return err
	}

	if float64(Balance.Current) < Order.Sum {
		logger.Log.Error("Not enough points")
		return balance.ErrNotEnoughPoints
	}

	if err := b.orderRepo.InsertOrder(ctx, Order); err != nil {
		return err
	}

	return nil
}

func (b BalanceUseCase) GetUserWithdrawals(ctx context.Context, uuid string) ([]models.Drawall, error) {
	return b.orderRepo.GetDrawalsByUUID(ctx, uuid)
}
