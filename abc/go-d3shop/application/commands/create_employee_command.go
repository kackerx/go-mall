package commands

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/domain/aggregates/salary"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// CreateEmployeeCommand 创建员工命令
type CreateEmployeeCommand struct {
	Name         string `validate:"required"`
	Email        string `validate:"required,email"`
	DepartmentID int64  `validate:"required,min=1"`
	Position     string `validate:"required"`
	BaseSalary   int    `validate:"required,min=1"`
}

// 确保实现ICommand接口
var _ mediator.ICommand = (*CreateEmployeeCommand)(nil)

// CreateEmployeeCommandHandler 创建员工命令处理器
type CreateEmployeeCommandHandler struct {
	employeeRepo   repositories.IEmployeeRepository
	departmentRepo repositories.IDepartmentRepository
}

// NewCreateEmployeeCommandHandler 创建命令处理器
func NewCreateEmployeeCommandHandler(
	employeeRepo repositories.IEmployeeRepository,
	departmentRepo repositories.IDepartmentRepository,
) *CreateEmployeeCommandHandler {
	return &CreateEmployeeCommandHandler{
		employeeRepo:   employeeRepo,
		departmentRepo: departmentRepo,
	}
}

// Handle 处理命令
func (h *CreateEmployeeCommandHandler) Handle(ctx context.Context, cmd CreateEmployeeCommand) (employee.EmployeeID, error) {
	// 验证部门是否存在
	dept, err := h.departmentRepo.GetByID(ctx, employee.NewDepartmentID(cmd.DepartmentID))
	if err != nil {
		return employee.EmployeeID{}, err
	}
	if dept == nil {
		return employee.EmployeeID{}, mediator.ErrInvalidRequest
	}

	// 创建员工聚合根
	emp, err := employee.NewEmployee(
		cmd.Name,
		cmd.Email,
		employee.NewDepartmentID(cmd.DepartmentID),
		cmd.Position,
		cmd.BaseSalary,
	)
	if err != nil {
		return employee.EmployeeID{}, err
	}

	// 保存到仓储
	err = h.employeeRepo.Add(ctx, emp)
	if err != nil {
		return employee.EmployeeID{}, err
	}

	return emp.ID, nil
}

// RegisterEmployeeCommandHandlers 注册员工相关的命令处理器到MediatorV2
func RegisterEmployeeCommandHandlers(
	med *mediator.MediatorV2,
	employeeRepo repositories.IEmployeeRepository,
	departmentRepo repositories.IDepartmentRepository,
	salaryRepo repositories.ISalaryRepository,
) {
	// 注册创建员工命令处理器
	med.RegisterCommandHandler("CreateEmployee", func(ctx context.Context, request interface{}) (interface{}, error) {
		cmd, ok := request.(*CreateEmployeeCommand)
		if !ok {
			return nil, mediator.ErrInvalidRequest
		}

		handler := NewCreateEmployeeCommandHandler(employeeRepo, departmentRepo)
		return handler.Handle(ctx, *cmd)
	})

	// 注册更新部门员工数命令处理器
	med.RegisterCommandHandler("UpdateDepartmentEmployeeCount", func(ctx context.Context, request interface{}) (interface{}, error) {
		cmd, ok := request.(*UpdateDepartmentEmployeeCountCommand)
		if !ok {
			return nil, mediator.ErrInvalidRequest
		}

		dept, err := departmentRepo.GetByID(ctx, cmd.DepartmentID)
		if err != nil {
			return nil, err
		}
		if dept != nil {
			if cmd.Delta > 0 {
				dept.AddEmployee()
			} else {
				dept.RemoveEmployee()
			}
			err = departmentRepo.Update(ctx, dept)
		}
		return nil, err
	})

	// 注册创建工资记录命令处理器
	med.RegisterCommandHandler("CreateSalaryRecord", func(ctx context.Context, request interface{}) (interface{}, error) {
		cmd, ok := request.(*CreateSalaryRecordCommand)
		if !ok {
			return nil, mediator.ErrInvalidRequest
		}

		record := salary.NewSalaryRecord(cmd.EmployeeID, cmd.BaseSalary)
		err := salaryRepo.Add(ctx, record)
		return record.ID, err
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
