package ddd

import (
	"time"
)

// IDomainEvent 领域事件接口
type IDomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

// IEntity 实体接口
type IEntity interface {
	GetID() interface{}
	GetDomainEvents() []IDomainEvent
	AddDomainEvent(event IDomainEvent)
	ClearDomainEvents()
}

// IAggregateRoot 聚合根接口
type IAggregateRoot interface {
	IEntity
}

// BaseEntity 基础实体
type BaseEntity struct {
	domainEvents []IDomainEvent
}

// GetDomainEvents 获取领域事件
func (e *BaseEntity) GetDomainEvents() []IDomainEvent {
	return e.domainEvents
}

// AddDomainEvent 添加领域事件
func (e *BaseEntity) AddDomainEvent(event IDomainEvent) {
	e.domainEvents = append(e.domainEvents, event)
}

// ClearDomainEvents 清除领域事件
func (e *BaseEntity) ClearDomainEvents() {
	e.domainEvents = []IDomainEvent{}
}

// BaseDomainEvent 基础领域事件
type BaseDomainEvent struct {
	occurredAt time.Time
}

// OccurredAt 事件发生时间
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// NewBaseDomainEvent 创建基础领域事件
func NewBaseDomainEvent() BaseDomainEvent {
	return BaseDomainEvent{
		occurredAt: time.Now(),
	}
}
