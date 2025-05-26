package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/order"
)

// IOrderRepository 订单仓储接口
type IOrderRepository interface {
	GetByID(ctx context.Context, id order.OrderID) (*order.Order, error)
	Add(ctx context.Context, order *order.Order) error
	Update(ctx context.Context, order *order.Order) error
	Delete(ctx context.Context, id order.OrderID) error
}
