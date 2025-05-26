package domain_event_handlers

import (
	"context"
	"log"

	"github.com/yourusername/go-d3shop/application/commands"
	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/domain/aggregates/salary"
	"github.com/yourusername/go-d3shop/domain/events"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// EmployeeJoinedDomainEventHandler 员工入职领域事件处理器
type EmployeeJoinedDomainEventHandler struct {
	mediator       *mediator.MediatorV2
	departmentRepo repositories.IDepartmentRepository
	salaryRepo     repositories.ISalaryRepository
}

// NewEmployeeJoinedDomainEventHandler 创建事件处理器
func NewEmployeeJoinedDomainEventHandler(
	mediator *mediator.MediatorV2,
	departmentRepo repositories.IDepartmentRepository,
	salaryRepo repositories.ISalaryRepository,
) *EmployeeJoinedDomainEventHandler {
	return &EmployeeJoinedDomainEventHandler{
		mediator:       mediator,
		departmentRepo: departmentRepo,
		salaryRepo:     salaryRepo,
	}
}

// Handle 处理事件 - 直接处理风格
func (h *EmployeeJoinedDomainEventHandler) Handle(ctx context.Context, event *events.EmployeeJoinedDomainEvent) error {
	// 从事件中获取员工信息
	emp, ok := event.Employee.(*employee.Employee)
	if !ok {
		return nil
	}

	log.Printf("处理员工入职领域事件: 员工ID=%d, 姓名=%s", emp.ID.Value(), emp.Name)

	// 1. 更新部门员工数
	dept, err := h.departmentRepo.GetByID(ctx, emp.DepartmentID)
	if err != nil {
		log.Printf("获取部门失败: %v", err)
		return err
	}
	if dept != nil {
		dept.AddEmployee()
		if err := h.departmentRepo.Update(ctx, dept); err != nil {
			log.Printf("更新部门员工数失败: %v", err)
			return err
		}
		log.Printf("部门[%s]员工数已更新: %d", dept.Name, dept.EmployeeCount)
	}

	// 2. 创建默认工资记录
	salaryRecord := salary.NewSalaryRecord(emp.ID, emp.BaseSalary)
	if err := h.salaryRepo.Add(ctx, salaryRecord); err != nil {
		log.Printf("创建工资记录失败: %v", err)
		return err
	}
	log.Printf("工资记录已创建: 员工ID=%d, 基本工资=%d", emp.ID.Value(), emp.BaseSalary)

	return nil
}

// HandleCommandStyle 处理事件 - 事件转命令风格
func (h *EmployeeJoinedDomainEventHandler) HandleCommandStyle(ctx context.Context, event *events.EmployeeJoinedDomainEvent) error {
	// 从事件中获取员工信息
	emp, ok := event.Employee.(*employee.Employee)
	if !ok {
		return nil
	}

	log.Printf("处理员工入职领域事件（命令风格）: 员工ID=%d", emp.ID.Value())

	// 1. 发送更新部门员工数命令
	_, err := h.mediator.Send(ctx, "UpdateDepartmentEmployeeCount", &commands.UpdateDepartmentEmployeeCountCommand{
		DepartmentID: emp.DepartmentID,
		Delta:        1,
	})
	if err != nil {
		log.Printf("更新部门员工数失败: %v", err)
		return err
	}

	// 2. 发送创建工资记录命令
	_, err = h.mediator.Send(ctx, "CreateSalaryRecord", &commands.CreateSalaryRecordCommand{
		EmployeeID: emp.ID,
		BaseSalary: emp.BaseSalary,
	})
	if err != nil {
		log.Printf("创建工资记录失败: %v", err)
		return err
	}

	return nil
}

// RegisterEmployeeDomainEventHandlers 注册员工相关的领域事件处理器
func RegisterEmployeeDomainEventHandlers(
	med *mediator.MediatorV2,
	departmentRepo repositories.IDepartmentRepository,
	salaryRepo repositories.ISalaryRepository,
	useCommandStyle bool,
) {
	handler := NewEmployeeJoinedDomainEventHandler(med, departmentRepo, salaryRepo)

	if useCommandStyle {
		// 事件转命令风格
		med.RegisterNotificationHandler("EmployeeJoined", func(ctx context.Context, notification interface{}) error {
			event, ok := notification.(*events.EmployeeJoinedDomainEvent)
			if !ok {
				return nil
			}
			return handler.HandleCommandStyle(ctx, event)
		})
	} else {
		// 直接处理风格
		med.RegisterNotificationHandler("EmployeeJoined", func(ctx context.Context, notification interface{}) error {
			event, ok := notification.(*events.EmployeeJoinedDomainEvent)
			if !ok {
				return nil
			}
			return handler.Handle(ctx, event)
		})
	}
}
