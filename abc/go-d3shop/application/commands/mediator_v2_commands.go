package commands

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// CreateOrderCommandV3 创建订单命令（MediatorV2版本）
type CreateOrderCommandV3 struct {
	Name  string
	Price int
	Count int
}

// OrderPaidCommandV2 订单支付命令（MediatorV2版本）
type OrderPaidCommandV2 struct {
	OrderID order.OrderID
}

// RegisterCommandHandlers 注册所有命令处理器到MediatorV2
func RegisterCommandHandlers(med *mediator.MediatorV2, orderRepo repositories.IOrderRepository) {
	// 注册创建订单命令处理器
	med.RegisterCommandHandler("CreateOrder", func(ctx context.Context, request interface{}) (interface{}, error) {
		cmd, ok := request.(*CreateOrderCommandV3)
		if !ok {
			return nil, mediator.ErrInvalidRequest
		}

		// 创建订单聚合根
		orderAgg := order.NewOrder(cmd.Name, cmd.Count)

		// 保存到仓储
		err := orderRepo.Add(ctx, orderAgg)
		if err != nil {
			return nil, err
		}

		// 发布领域事件（这里由仓储或服务处理）

		return orderAgg.ID, nil
	})

	// 注册订单支付命令处理器
	med.RegisterCommandHandler("PayOrder", func(ctx context.Context, request interface{}) (interface{}, error) {
		cmd, ok := request.(*OrderPaidCommandV2)
		if !ok {
			return nil, mediator.ErrInvalidRequest
		}

		// 获取订单
		orderAgg, err := orderRepo.GetByID(ctx, cmd.OrderID)
		if err != nil {
			return nil, err
		}

		// 执行支付
		err = orderAgg.OrderPaid()
		if err != nil {
			return nil, err
		}

		// 更新订单
		err = orderRepo.Update(ctx, orderAgg)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
}
