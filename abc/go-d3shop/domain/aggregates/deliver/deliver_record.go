package deliver

import (
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/pkg/ddd"
)

// DeliverRecordID 发货记录ID
type DeliverRecordID struct {
	ddd.Int64StronglyTypedId
}

// NewDeliverRecordID 创建发货记录ID
func NewDeliverRecordID(id int64) DeliverRecordID {
	return DeliverRecordID{ddd.NewInt64StronglyTypedId(id)}
}

// DeliverRecord 发货记录聚合根
type DeliverRecord struct {
	ddd.BaseEntity
	ID      DeliverRecordID `gorm:"primaryKey;column:id"`
	OrderID order.OrderID   `gorm:"column:order_id"`
}

// NewDeliverRecord 创建新的发货记录
func NewDeliverRecord(orderID order.OrderID) *DeliverRecord {
	return &DeliverRecord{
		OrderID: orderID,
	}
}

// GetID 获取ID
func (d *DeliverRecord) GetID() interface{} {
	return d.ID
}

// TableName 指定表名
func (DeliverRecord) TableName() string {
	return "deliver_records"
}
