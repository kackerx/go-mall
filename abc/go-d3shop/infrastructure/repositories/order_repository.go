package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/infrastructure/persistence"
)

// OrderRepository 订单仓储实现
type OrderRepository struct {
	dbContext *persistence.DbContext
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(dbContext *persistence.DbContext) repositories.IOrderRepository {
	return &OrderRepository{
		dbContext: dbContext,
	}
}

// GetByID 根据ID获取订单
func (r *OrderRepository) GetByID(ctx context.Context, id order.OrderID) (*order.Order, error) {
	var orderEntity order.Order
	err := r.dbContext.DB().WithContext(ctx).First(&orderEntity, "id = ?", id.Value()).Error
	if err != nil {
		return nil, err
	}
	return &orderEntity, nil
}

// Add 添加订单
func (r *OrderRepository) Add(ctx context.Context, orderAgg *order.Order) error {
	// 保存订单
	err := r.dbContext.DB().WithContext(ctx).Create(orderAgg).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, orderAgg)
}

// Update 更新订单
func (r *OrderRepository) Update(ctx context.Context, orderAgg *order.Order) error {
	// 更新订单
	err := r.dbContext.DB().WithContext(ctx).Save(orderAgg).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, orderAgg)
}

// Delete 删除订单
func (r *OrderRepository) Delete(ctx context.Context, id order.OrderID) error {
	return r.dbContext.DB().WithContext(ctx).Delete(&order.Order{}, "id = ?", id.Value()).Error
}
