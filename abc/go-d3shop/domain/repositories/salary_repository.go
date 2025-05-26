package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/domain/aggregates/salary"
)

// ISalaryRepository 工资仓储接口
type ISalaryRepository interface {
	GetByID(ctx context.Context, id salary.SalaryRecordID) (*salary.SalaryRecord, error)
	GetByEmployeeID(ctx context.Context, employeeID employee.EmployeeID) ([]*salary.SalaryRecord, error)
	Add(ctx context.Context, record *salary.SalaryRecord) error
	Update(ctx context.Context, record *salary.SalaryRecord) error
}
