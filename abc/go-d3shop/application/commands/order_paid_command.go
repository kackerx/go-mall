package commands

import (
	"context"
	"errors"

	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// OrderPaidCommand 订单支付命令
type OrderPaidCommand struct {
	OrderID order.OrderID
}

// 确保实现ICommand接口
var _ mediator.ICommand = (*OrderPaidCommand)(nil)

// OrderPaidCommandHandler 订单支付命令处理器
type OrderPaidCommandHandler struct {
	orderRepo repositories.IOrderRepository
}

// NewOrderPaidCommandHandler 创建命令处理器
func NewOrderPaidCommandHandler(orderRepo repositories.IOrderRepository) *OrderPaidCommandHandler {
	return &OrderPaidCommandHandler{
		orderRepo: orderRepo,
	}
}

// Handle 处理命令
func (h *OrderPaidCommandHandler) Handle(ctx context.Context, cmd OrderPaidCommand) error {
	// 获取订单
	orderAgg, err := h.orderRepo.GetByID(ctx, cmd.OrderID)
	if err != nil {
		return err
	}
	if orderAgg == nil {
		return errors.New("订单不存在")
	}

	// 执行支付
	err = orderAgg.OrderPaid()
	if err != nil {
		return err
	}

	// 更新订单
	return h.orderRepo.Update(ctx, orderAgg)
}
