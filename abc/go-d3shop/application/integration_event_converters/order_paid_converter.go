package integration_event_converters

import (
	"github.com/yourusername/go-d3shop/application/integration_events"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/events"
)

// OrderPaidIntegrationEventConverter 订单支付事件转换器
type OrderPaidIntegrationEventConverter struct{}

// NewOrderPaidIntegrationEventConverter 创建转换器
func NewOrderPaidIntegrationEventConverter() *OrderPaidIntegrationEventConverter {
	return &OrderPaidIntegrationEventConverter{}
}

// Convert 转换领域事件为集成事件
func (c *OrderPaidIntegrationEventConverter) Convert(domainEvent *events.OrderPaidDomainEvent) *integration_events.OrderPaidIntegrationEvent {
	// 从领域事件中提取订单信息
	orderAgg, ok := domainEvent.Order.(*order.Order)
	if !ok {
		return nil
	}

	return integration_events.NewOrderPaidIntegrationEvent(orderAgg.ID)
}
