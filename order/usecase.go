package order

import (
	"context"

	"github.com/0loff/grade_gophermart/models"
)

type UseCase interface {
	SetOrder(ctx context.Context, order, uuid string) error
	GetUserOrders(ctx context.Context, uuid string) ([]models.Order, error)
}
