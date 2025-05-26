package integration_events

import (
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/pkg/events"
)

// EmployeeJoinedIntegrationEvent 员工入职集成事件
// 用于通知其他服务（如通知服务）
type EmployeeJoinedIntegrationEvent struct {
	events.BaseIntegrationEvent
	EmployeeID     employee.EmployeeID   `json:"employeeId"`
	EmployeeName   string                `json:"employeeName"`
	EmployeeEmail  string                `json:"employeeEmail"`
	DepartmentID   employee.DepartmentID `json:"departmentId"`
	DepartmentName string                `json:"departmentName"`
	ManagerID      employee.EmployeeID   `json:"managerId"`
	ManagerEmail   string                `json:"managerEmail"`
	Position       string                `json:"position"`
}

// NewEmployeeJoinedIntegrationEvent 创建员工入职集成事件
func NewEmployeeJoinedIntegrationEvent(
	employeeID employee.EmployeeID,
	employeeName, employeeEmail string,
	departmentID employee.DepartmentID,
	departmentName string,
	managerID employee.EmployeeID,
	managerEmail string,
	position string,
) *EmployeeJoinedIntegrationEvent {
	return &EmployeeJoinedIntegrationEvent{
		BaseIntegrationEvent: events.BaseIntegrationEvent{
			ID:         uuid.New().String(),
			OccurredOn: time.Now(),
		},
		EmployeeID:     employeeID,
		EmployeeName:   employeeName,
		EmployeeEmail:  employeeEmail,
		DepartmentID:   departmentID,
		DepartmentName: departmentName,
		ManagerID:      managerID,
		ManagerEmail:   managerEmail,
		Position:       position,
	}
}

// EventName 事件名称
func (e *EmployeeJoinedIntegrationEvent) EventName() string {
	return "EmployeeJoinedIntegrationEvent"
}
