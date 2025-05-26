package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/department"
	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
)

// IDepartmentRepository 部门仓储接口
type IDepartmentRepository interface {
	GetByID(ctx context.Context, id employee.DepartmentID) (*department.Department, error)
	Add(ctx context.Context, dept *department.Department) error
	Update(ctx context.Context, dept *department.Department) error
}
