package services

import (
	"context"
	"errors"
	"log"

	"github.com/yourusername/go-d3shop/application/integration_events"
	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/domain/aggregates/salary"
	"github.com/yourusername/go-d3shop/domain/events"
	"github.com/yourusername/go-d3shop/domain/repositories"
	pkgEvents "github.com/yourusername/go-d3shop/pkg/events"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// EmployeeService 员工服务
type EmployeeService struct {
	employeeRepo   repositories.IEmployeeRepository
	departmentRepo repositories.IDepartmentRepository
	salaryRepo     repositories.ISalaryRepository
	mediator       *mediator.MediatorV2
	eventPublisher pkgEvents.IIntegrationEventPublisher
}

// NewEmployeeService 创建员工服务
func NewEmployeeService(
	employeeRepo repositories.IEmployeeRepository,
	departmentRepo repositories.IDepartmentRepository,
	salaryRepo repositories.ISalaryRepository,
	mediator *mediator.MediatorV2,
	eventPublisher pkgEvents.IIntegrationEventPublisher,
) *EmployeeService {
	return &EmployeeService{
		employeeRepo:   employeeRepo,
		departmentRepo: departmentRepo,
		salaryRepo:     salaryRepo,
		mediator:       mediator,
		eventPublisher: eventPublisher,
	}
}

// CreateEmployeeRequest 创建员工请求
type CreateEmployeeRequest struct {
	Name         string
	Email        string
	DepartmentID int64
	Position     string
	BaseSalary   int
}

// CreateEmployee 创建员工
func (s *EmployeeService) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (employee.EmployeeID, error) {
	// 验证部门是否存在
	dept, err := s.departmentRepo.GetByID(ctx, employee.NewDepartmentID(req.DepartmentID))
	if err != nil {
		return employee.EmployeeID{}, err
	}
	if dept == nil {
		return employee.EmployeeID{}, errors.New("部门不存在")
	}

	// 创建员工聚合根
	emp, err := employee.NewEmployee(
		req.Name,
		req.Email,
		employee.NewDepartmentID(req.DepartmentID),
		req.Position,
		req.BaseSalary,
	)
	if err != nil {
		return employee.EmployeeID{}, err
	}

	// 保存员工
	if err := s.employeeRepo.Add(ctx, emp); err != nil {
		return employee.EmployeeID{}, err
	}

	// 发布领域事件
	for _, event := range emp.GetDomainEvents() {
		if domainEvent, ok := event.(interface{ EventName() string }); ok {
			err := s.mediator.Publish(ctx, domainEvent.EventName(), event)
			if err != nil {
				log.Printf("发布领域事件失败: %v", err)
			}
		}
	}
	emp.ClearDomainEvents()

	// 发布集成事件（通知其他服务）
	// 获取部门经理信息用于通知
	var managerEmail string
	if manager, err := s.employeeRepo.GetByID(ctx, dept.ManagerID); err == nil && manager != nil {
		managerEmail = manager.Email
	}

	integrationEvent := integration_events.NewEmployeeJoinedIntegrationEvent(
		emp.ID,
		emp.Name,
		emp.Email,
		dept.ID,
		dept.Name,
		dept.ManagerID,
		managerEmail,
		emp.Position,
	)

	// 异步发布集成事件
	if err := s.eventPublisher.PublishAsync(ctx, integrationEvent); err != nil {
		log.Printf("发布集成事件失败: %v", err)
	}

	return emp.ID, nil
}

// RegisterEventHandlers 注册事件处理器
func (s *EmployeeService) RegisterEventHandlers() {
	// 处理员工入职领域事件 - 更新部门员工数
	s.mediator.RegisterNotificationHandler("EmployeeJoined", func(ctx context.Context, notification interface{}) error {
		event, ok := notification.(*events.EmployeeJoinedDomainEvent)
		if !ok {
			return nil
		}

		emp, ok := event.Employee.(*employee.Employee)
		if !ok {
			return nil
		}

		// 领域事件1：更新部门员工数
		dept, err := s.departmentRepo.GetByID(ctx, emp.DepartmentID)
		if err != nil {
			return err
		}
		if dept != nil {
			dept.AddEmployee()
			if err := s.departmentRepo.Update(ctx, dept); err != nil {
				log.Printf("更新部门员工数失败: %v", err)
			}
		}

		// 领域事件2：创建默认工资记录
		salaryRecord := salary.NewSalaryRecord(emp.ID, emp.BaseSalary)
		if err := s.salaryRepo.Add(ctx, salaryRecord); err != nil {
			log.Printf("创建工资记录失败: %v", err)
		}

		return nil
	})
}

// 也可以将领域事件处理分开，使用事件转命令模式
func (s *EmployeeService) RegisterEventHandlersCommandStyle() {
	// 处理员工入职 - 更新部门
	s.mediator.RegisterNotificationHandler("EmployeeJoined", func(ctx context.Context, notification interface{}) error {
		event, ok := notification.(*events.EmployeeJoinedDomainEvent)
		if !ok {
			return nil
		}

		emp, ok := event.Employee.(*employee.Employee)
		if !ok {
			return nil
		}

		// 发送更新部门命令
		_, err := s.mediator.Send(ctx, "UpdateDepartmentEmployeeCount", &UpdateDepartmentEmployeeCountCommand{
			DepartmentID: emp.DepartmentID,
			Delta:        1,
		})
		return err
	})

	// 处理员工入职 - 创建工资记录
	s.mediator.RegisterNotificationHandler("EmployeeJoined", func(ctx context.Context, notification interface{}) error {
		event, ok := notification.(*events.EmployeeJoinedDomainEvent)
		if !ok {
			return nil
		}

		emp, ok := event.Employee.(*employee.Employee)
		if !ok {
			return nil
		}

		// 发送创建工资记录命令
		_, err := s.mediator.Send(ctx, "CreateSalaryRecord", &CreateSalaryRecordCommand{
			EmployeeID: emp.ID,
			BaseSalary: emp.BaseSalary,
		})
		return err
	})
}

// UpdateDepartmentEmployeeCountCommand 更新部门员工数命令
type UpdateDepartmentEmployeeCountCommand struct {
	DepartmentID employee.DepartmentID
	Delta        int // +1 或 -1
}

// CreateSalaryRecordCommand 创建工资记录命令
type CreateSalaryRecordCommand struct {
	EmployeeID employee.EmployeeID
	BaseSalary int
}
