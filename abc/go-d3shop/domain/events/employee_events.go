package events

import (
	"github.com/yourusername/go-d3shop/pkg/ddd"
)

// EmployeeJoinedDomainEvent 员工入职领域事件
type EmployeeJoinedDomainEvent struct {
	ddd.BaseDomainEvent
	Employee interface{} // 避免循环依赖，使用interface{}
}

// NewEmployeeJoinedDomainEvent 创建员工入职事件
func NewEmployeeJoinedDomainEvent(employee interface{}) *EmployeeJoinedDomainEvent {
	return &EmployeeJoinedDomainEvent{
		BaseDomainEvent: ddd.NewBaseDomainEvent(),
		Employee:        employee,
	}
}

// EventName 事件名称
func (e *EmployeeJoinedDomainEvent) EventName() string {
	return "EmployeeJoined"
}
