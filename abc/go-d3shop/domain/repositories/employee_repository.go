package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
)

// IEmployeeRepository 员工仓储接口
type IEmployeeRepository interface {
	GetByID(ctx context.Context, id employee.EmployeeID) (*employee.Employee, error)
	GetByEmail(ctx context.Context, email string) (*employee.Employee, error)
	Add(ctx context.Context, emp *employee.Employee) error
	Update(ctx context.Context, emp *employee.Employee) error
	Delete(ctx context.Context, id employee.EmployeeID) error
}
