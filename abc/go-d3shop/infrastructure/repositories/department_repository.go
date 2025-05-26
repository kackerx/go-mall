package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/department"
	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/infrastructure/persistence"
)

// DepartmentRepository 部门仓储实现
type DepartmentRepository struct {
	dbContext *persistence.DbContext
}

// NewDepartmentRepository 创建部门仓储
func NewDepartmentRepository(dbContext *persistence.DbContext) repositories.IDepartmentRepository {
	return &DepartmentRepository{
		dbContext: dbContext,
	}
}

// GetByID 根据ID获取部门
func (r *DepartmentRepository) GetByID(ctx context.Context, id employee.DepartmentID) (*department.Department, error) {
	var dept department.Department
	err := r.dbContext.DB().WithContext(ctx).First(&dept, "id = ?", id.Value()).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

// Add 添加部门
func (r *DepartmentRepository) Add(ctx context.Context, dept *department.Department) error {
	// 保存部门
	err := r.dbContext.DB().WithContext(ctx).Create(dept).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, dept)
}

// Update 更新部门
func (r *DepartmentRepository) Update(ctx context.Context, dept *department.Department) error {
	// 更新部门
	err := r.dbContext.DB().WithContext(ctx).Save(dept).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, dept)
}
