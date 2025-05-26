package persistence

import (
	"context"

	"github.com/yourusername/go-d3shop/pkg/ddd"
	"github.com/yourusername/go-d3shop/pkg/mediator"
	"gorm.io/gorm"
)

// DbContext 数据库上下文
type DbContext struct {
	db       *gorm.DB
	mediator mediator.IMediator
}

// NewDbContext 创建数据库上下文
func NewDbContext(db *gorm.DB, mediator mediator.IMediator) *DbContext {
	return &DbContext{
		db:       db,
		mediator: mediator,
	}
}

// DB 获取数据库连接
func (c *DbContext) DB() *gorm.DB {
	return c.db
}

// SaveChangesAsync 保存更改并发布领域事件
func (c *DbContext) SaveChangesAsync(ctx context.Context) error {
	// 开启事务
	tx := c.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 获取所有待发布的领域事件
	var domainEvents []ddd.IDomainEvent

	// 这里需要通过某种方式收集所有实体的领域事件
	// 在实际实现中，可能需要追踪所有被修改的实体

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// 发布领域事件（在事务提交后）
	for _, event := range domainEvents {
		if err := c.mediator.Publish(ctx, event); err != nil {
			// 记录错误但不影响主流程
			// 在实际应用中应该有更好的错误处理机制
			continue
		}
	}

	return nil
}

// PublishDomainEvents 发布实体的领域事件
func (c *DbContext) PublishDomainEvents(ctx context.Context, entity ddd.IEntity) error {
	events := entity.GetDomainEvents()
	for _, event := range events {
		if err := c.mediator.Publish(ctx, event); err != nil {
			return err
		}
	}
	entity.ClearDomainEvents()
	return nil
}
