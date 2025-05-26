package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/deliver"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/infrastructure/persistence"
)

// DeliverRecordRepository 发货记录仓储实现
type DeliverRecordRepository struct {
	dbContext *persistence.DbContext
}

// NewDeliverRecordRepository 创建发货记录仓储
func NewDeliverRecordRepository(dbContext *persistence.DbContext) repositories.IDeliverRecordRepository {
	return &DeliverRecordRepository{
		dbContext: dbContext,
	}
}

// GetByID 根据ID获取发货记录
func (r *DeliverRecordRepository) GetByID(ctx context.Context, id deliver.DeliverRecordID) (*deliver.DeliverRecord, error) {
	var record deliver.DeliverRecord
	err := r.dbContext.DB().WithContext(ctx).First(&record, "id = ?", id.Value()).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// Add 添加发货记录
func (r *DeliverRecordRepository) Add(ctx context.Context, record *deliver.DeliverRecord) error {
	// 保存发货记录
	err := r.dbContext.DB().WithContext(ctx).Create(record).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, record)
}

// Update 更新发货记录
func (r *DeliverRecordRepository) Update(ctx context.Context, record *deliver.DeliverRecord) error {
	// 更新发货记录
	err := r.dbContext.DB().WithContext(ctx).Save(record).Error
	if err != nil {
		return err
	}

	// 发布领域事件
	return r.dbContext.PublishDomainEvents(ctx, record)
}
