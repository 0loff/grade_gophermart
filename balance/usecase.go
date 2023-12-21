package balance

import (
	"context"

	"github.com/0loff/grade_gophermart/models"
)

type UseCase interface {
	GetUserBalance(ctx context.Context, uuid string) (models.Balance, error)
	SetOrderWithdraw(ctx context.Context, Order models.Order) error
	GetUserWithdrawals(ctx context.Context, uuid string) ([]models.Drawall, error)
}
