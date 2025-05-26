package order

import (
	"errors"

	"github.com/yourusername/go-d3shop/domain/events"
	"github.com/yourusername/go-d3shop/pkg/ddd"
)

// OrderID 订单ID
type OrderID struct {
	ddd.Int64StronglyTypedId
}

// NewOrderID 创建订单ID
func NewOrderID(id int64) OrderID {
	return OrderID{ddd.NewInt64StronglyTypedId(id)}
}

// Order 订单聚合根
type Order struct {
	ddd.BaseEntity
	ID    OrderID `gorm:"primaryKey;column:id"`
	Name  string  `gorm:"column:name"`
	Count int     `gorm:"column:count"`
	Paid  bool    `gorm:"column:paid"`
}

// NewOrder 创建新订单
func NewOrder(name string, count int) *Order {
	order := &Order{
		Name:  name,
		Count: count,
		Paid:  false,
	}

	// 添加订单创建领域事件
	order.AddDomainEvent(events.NewOrderCreatedDomainEvent(order))

	return order
}

// GetID 获取订单ID
func (o *Order) GetID() interface{} {
	return o.ID
}

// OrderPaid 订单支付
func (o *Order) OrderPaid() error {
	if o.Paid {
		return errors.New("订单已经支付")
	}

	o.Paid = true

	// 添加订单支付领域事件
	o.AddDomainEvent(events.NewOrderPaidDomainEvent(o))

	return nil
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}
