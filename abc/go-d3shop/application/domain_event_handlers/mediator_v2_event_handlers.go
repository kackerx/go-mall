package domain_event_handlers

import (
	"context"

	"github.com/yourusername/go-d3shop/application/commands"
	"github.com/yourusername/go-d3shop/domain/aggregates/deliver"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/events"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// RegisterDomainEventHandlers 注册所有领域事件处理器到MediatorV2
func RegisterDomainEventHandlers(
	med *mediator.MediatorV2,
	orderRepo repositories.IOrderRepository,
	deliverRepo repositories.IDeliverRecordRepository,
) {
	// 方式1：事件转命令（保持与原MediatR一致的风格）
	RegisterOrderCreatedHandler_CommandStyle(med)

	// 方式2：直接处理（更Go风格）
	RegisterOrderPaidHandler_DirectStyle(med, orderRepo)
}

// RegisterOrderCreatedHandler_CommandStyle 注册订单创建事件处理器（命令风格）
// 这种方式保持了事件与命令的解耦，适合复杂的业务流程
func RegisterOrderCreatedHandler_CommandStyle(med *mediator.MediatorV2) {
	med.RegisterNotificationHandler("OrderCreated", func(ctx context.Context, notification interface{}) error {
		event, ok := notification.(*events.OrderCreatedDomainEvent)
		if !ok {
			return nil
		}

		// 从事件中获取订单信息
		orderAgg, ok := event.Order.(*order.Order)
		if !ok {
			return nil
		}

		// 转换为命令并发送
		cmd := &commands.DeliverGoodsCommand{
			OrderID: orderAgg.ID,
		}

		// 通过MediatorV2发送命令
		// 注意：这里需要先注册DeliverGoods命令处理器
		_, err := med.Send(ctx, "DeliverGoods", cmd)
		return err
	})
}

// RegisterOrderPaidHandler_DirectStyle 注册订单支付事件处理器（直接处理风格）
// 这种方式更直接，适合简单的业务逻辑
func RegisterOrderPaidHandler_DirectStyle(med *mediator.MediatorV2, orderRepo repositories.IOrderRepository) {
	med.RegisterNotificationHandler("OrderPaid", func(ctx context.Context, notification interface{}) error {
		event, ok := notification.(*events.OrderPaidDomainEvent)
		if !ok {
			return nil
		}

		// 直接执行业务逻辑
		// 例如：更新订单统计、发送通知等
		orderAgg, ok := event.Order.(*order.Order)
		if !ok {
			return nil
		}

		// 这里可以直接调用其他服务或仓储
		// 比如：
		// - 更新用户积分
		// - 发送订单支付成功通知
		// - 更新销售统计

		// 示例：记录日志
		println("订单支付成功，订单ID:", orderAgg.ID.Value())

		return nil
	})
}

// RegisterDeliverGoodsCommandHandler 注册发货命令处理器
// 配合命令风格的事件处理器使用
func RegisterDeliverGoodsCommandHandler(med *mediator.MediatorV2, deliverRepo repositories.IDeliverRecordRepository) {
	med.RegisterCommandHandler("DeliverGoods", func(ctx context.Context, request interface{}) (interface{}, error) {
		cmd, ok := request.(*commands.DeliverGoodsCommand)
		if !ok {
			return nil, mediator.ErrInvalidRequest
		}

		// 创建发货记录
		record := deliver.NewDeliverRecord(cmd.OrderID)

		// 保存到仓储
		err := deliverRepo.Add(ctx, record)
		if err != nil {
			return nil, err
		}

		return record.ID, nil
	})
}

// 混合使用示例：
// 1. 简单的业务逻辑（如创建关联记录）可以直接在事件处理器中完成
// 2. 复杂的业务流程（如涉及多个聚合根的操作）可以转换为命令
// 3. 需要事务保证的操作应该在同一个事件处理器中完成
// 4. 可以异步处理的操作可以发布集成事件

// OrderCreatedEventHandlerV2 展示更复杂的事件处理
type OrderCreatedEventHandlerV2 struct {
	deliverRepo repositories.IDeliverRecordRepository
	mediator    *mediator.MediatorV2
}

func (h *OrderCreatedEventHandlerV2) Handle(ctx context.Context, notification interface{}) error {
	event, ok := notification.(*events.OrderCreatedDomainEvent)
	if !ok {
		return nil
	}

	orderAgg, ok := event.Order.(*order.Order)
	if !ok {
		return nil
	}

	// 1. 直接创建发货记录（简单逻辑）
	record := deliver.NewDeliverRecord(orderAgg.ID)
	if err := h.deliverRepo.Add(ctx, record); err != nil {
		return err
	}

	// 2. 如果需要触发其他复杂流程，可以发送命令
	// 例如：初始化库存预留
	// _, err := h.mediator.Send(ctx, "ReserveInventory", &ReserveInventoryCommand{OrderID: orderAgg.ID})

	return nil
}
