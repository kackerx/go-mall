package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/infrastructure/persistence"
)

// EmployeeRepository 员工仓储实现
type EmployeeRepository struct {
	dbContext *persistence.DbContext
}

// NewEmployeeRepository 创建员工仓储
func NewEmployeeRepository(dbContext *persistence.DbContext) repositories.IEmployeeRepository {
	return &EmployeeRepository{
		dbContext: dbContext,
	}
}

// GetByID 根据ID获取员工
func (r *EmployeeRepository) GetByID(ctx context.Context, id employee.EmployeeID) (*employee.Employee, error) {
	var emp employee.Employee
	err := r.dbContext.DB().WithContext(ctx).First(&emp, "id = ?", id.Value()).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// GetByEmail 根据邮箱获取员工
func (r *EmployeeRepository) GetByEmail(ctx context.Context, email string) (*employee.Employee, error) {
	var emp employee.Employee
	err := r.dbContext.DB().WithContext(ctx).Where("email = ?", email).First(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

// Add 添加员工
func (r *EmployeeRepository) Add(ctx context.Context, emp *employee.Employee) error {
	// 保存员工
	err := r.dbContext.DB().WithContext(ctx).Create(emp).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, emp)
}

// Update 更新员工
func (r *EmployeeRepository) Update(ctx context.Context, emp *employee.Employee) error {
	// 更新员工
	err := r.dbContext.DB().WithContext(ctx).Save(emp).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, emp)
}

// Delete 删除员工
func (r *EmployeeRepository) Delete(ctx context.Context, id employee.EmployeeID) error {
	return r.dbContext.DB().WithContext(ctx).Delete(&employee.Employee{}, "id = ?", id.Value()).Error
}
