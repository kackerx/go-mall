package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/domain/aggregates/salary"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/infrastructure/persistence"
)

// SalaryRepository 工资仓储实现
type SalaryRepository struct {
	dbContext *persistence.DbContext
}

// NewSalaryRepository 创建工资仓储
func NewSalaryRepository(dbContext *persistence.DbContext) repositories.ISalaryRepository {
	return &SalaryRepository{
		dbContext: dbContext,
	}
}

// GetByID 根据ID获取工资记录
func (r *SalaryRepository) GetByID(ctx context.Context, id salary.SalaryRecordID) (*salary.SalaryRecord, error) {
	var record salary.SalaryRecord
	err := r.dbContext.DB().WithContext(ctx).First(&record, "id = ?", id.Value()).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetByEmployeeID 根据员工ID获取工资记录
func (r *SalaryRepository) GetByEmployeeID(ctx context.Context, employeeID employee.EmployeeID) ([]*salary.SalaryRecord, error) {
	var records []*salary.SalaryRecord
	err := r.dbContext.DB().WithContext(ctx).Where("employee_id = ?", employeeID.Value()).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// Add 添加工资记录
func (r *SalaryRepository) Add(ctx context.Context, record *salary.SalaryRecord) error {
	// 保存工资记录
	err := r.dbContext.DB().WithContext(ctx).Create(record).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, record)
}

// Update 更新工资记录
func (r *SalaryRepository) Update(ctx context.Context, record *salary.SalaryRecord) error {
	// 更新工资记录
	err := r.dbContext.DB().WithContext(ctx).Save(record).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, record)
}
