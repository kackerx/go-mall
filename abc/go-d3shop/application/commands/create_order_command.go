package commands

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// CreateOrderCommand 创建订单命令
type CreateOrderCommand struct {
	Name  string
	Price int
	Count int
}

// 确保实现ICommand接口
var _ mediator.ICommand = (*CreateOrderCommand)(nil)

// CreateOrderCommandHandler 创建订单命令处理器
type CreateOrderCommandHandler struct {
	orderRepo repositories.IOrderRepository
}

// NewCreateOrderCommandHandler 创建命令处理器
func NewCreateOrderCommandHandler(orderRepo repositories.IOrderRepository) *CreateOrderCommandHandler {
	return &CreateOrderCommandHandler{
		orderRepo: orderRepo,
	}
}

// Handle 处理命令
func (h *CreateOrderCommandHandler) Handle(ctx context.Context, cmd CreateOrderCommand) (order.OrderID, error) {
	// 创建订单聚合根
	orderAgg := order.NewOrder(cmd.Name, cmd.Count)

	// 保存到仓储
	err := h.orderRepo.Add(ctx, orderAgg)
	if err != nil {
		return order.OrderID{}, err
	}

	return orderAgg.ID, nil
}
