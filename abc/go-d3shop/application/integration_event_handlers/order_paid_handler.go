package integration_event_handlers

import (
	"context"

	"github.com/yourusername/go-d3shop/application/commands"
	"github.com/yourusername/go-d3shop/application/integration_events"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// OrderPaidIntegrationEventHandler 订单支付集成事件处理器
type OrderPaidIntegrationEventHandler struct {
	mediator mediator.IMediator
}

// NewOrderPaidIntegrationEventHandler 创建事件处理器
func NewOrderPaidIntegrationEventHandler(mediator mediator.IMediator) *OrderPaidIntegrationEventHandler {
	return &OrderPaidIntegrationEventHandler{
		mediator: mediator,
	}
}

// HandleAsync 处理事件
func (h *OrderPaidIntegrationEventHandler) HandleAsync(ctx context.Context, event *integration_events.OrderPaidIntegrationEvent) error {
	// 发送订单支付命令
	cmd := commands.OrderPaidCommand{
		OrderID: event.OrderID,
	}

	_, err := h.mediator.Send(ctx, cmd)
	return err
}
