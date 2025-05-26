package events

import (
	"github.com/yourusername/go-d3shop/pkg/ddd"
)

// OrderCreatedDomainEvent 订单创建领域事件
type OrderCreatedDomainEvent struct {
	ddd.BaseDomainEvent
	Order interface{} // 避免循环依赖，使用interface{}
}

// NewOrderCreatedDomainEvent 创建订单创建事件
func NewOrderCreatedDomainEvent(order interface{}) *OrderCreatedDomainEvent {
	return &OrderCreatedDomainEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		Order:           order,
	}
}

// EventName 事件名称
func (e *OrderCreatedDomainEvent) EventName() string {
	return "OrderCreated"
}

// OrderPaidDomainEvent 订单支付领域事件
type OrderPaidDomainEvent struct {
	ddd.BaseDomainEvent
	Order interface{} // 避免循环依赖，使用interface{}
}

// NewOrderPaidDomainEvent 创建订单支付事件
func NewOrderPaidDomainEvent(order interface{}) *OrderPaidDomainEvent {
	return &OrderPaidDomainEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		Order:           order,
	}
}

// EventName 事件名称
func (e *OrderPaidDomainEvent) EventName() string {
	return "OrderPaid"
}
