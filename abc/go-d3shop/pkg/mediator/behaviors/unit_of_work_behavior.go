package behaviors

import (
	"context"

	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// IUnitOfWork 工作单元接口
type IUnitOfWork interface {
	BeginTransaction() error
	Commit() error
	Rollback() error
}

// UnitOfWorkBehavior 工作单元管道行为
type UnitOfWorkBehavior struct {
	unitOfWork IUnitOfWork
}

// NewUnitOfWorkBehavior 创建工作单元行为
func NewUnitOfWorkBehavior(unitOfWork IUnitOfWork) *UnitOfWorkBehavior {
	return &UnitOfWorkBehavior{
		unitOfWork: unitOfWork,
	}
}

// Handle 处理请求
func (b *UnitOfWorkBehavior) Handle(ctx context.Context, request mediator.IRequest, next mediator.RequestHandlerFunc) (interface{}, error) {
	// 开始事务
	err := b.unitOfWork.BeginTransaction()
	if err != nil {
		return nil, err
	}

	// 执行下一个处理器
	result, err := next(ctx)
	if err != nil {
		// 回滚事务
		b.unitOfWork.Rollback()
		return nil, err
	}

	// 提交事务
	err = b.unitOfWork.Commit()
	if err != nil {
		b.unitOfWork.Rollback()
		return nil, err
	}

	return result, nil
}
