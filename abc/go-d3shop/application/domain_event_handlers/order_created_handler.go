package domain_event_handlers

import (
	"context"

	"github.com/yourusername/go-d3shop/application/commands"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/events"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// OrderCreatedDomainEventHandler 订单创建领域事件处理器
type OrderCreatedDomainEventHandler struct {
	mediator mediator.IMediator
}

// NewOrderCreatedDomainEventHandler 创建事件处理器
func NewOrderCreatedDomainEventHandler(mediator mediator.IMediator) *OrderCreatedDomainEventHandler {
	return &OrderCreatedDomainEventHandler{
		mediator: mediator,
	}
}

// Handle 处理事件
func (h *OrderCreatedDomainEventHandler) Handle(ctx context.Context, event *events.OrderCreatedDomainEvent) error {
	// 从事件中获取订单信息
	orderAgg, ok := event.Order.(*order.Order)
	if !ok {
		return nil
	}

	// 发送发货命令
	cmd := commands.DeliverGoodsCommand{
		OrderID: orderAgg.ID,
	}

	_, err := h.mediator.Send(ctx, cmd)
	return err
}
