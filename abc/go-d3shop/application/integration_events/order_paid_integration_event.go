package integration_events

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/pkg/events"
)

// OrderPaidIntegrationEvent 订单支付集成事件
type OrderPaidIntegrationEvent struct {
	events.BaseIntegrationEvent
	OrderID order.OrderID `json:"orderId"`
}

// NewOrderPaidIntegrationEvent 创建订单支付集成事件
func NewOrderPaidIntegrationEvent(orderID order.OrderID) *OrderPaidIntegrationEvent {
	return &OrderPaidIntegrationEvent{
		BaseIntegrationEvent: events.BaseIntegrationEvent{
			ID:         uuid.New().String(),
			OccurredOn: time.Now(),
		},
		OrderID: orderID,
	}
}

// EventName 事件名称
func (e *OrderPaidIntegrationEvent) EventName() string {
	return "OrderPaidIntegrationEvent"
}
