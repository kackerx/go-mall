package repositories

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/deliver"
)

// IDeliverRecordRepository 发货记录仓储接口
type IDeliverRecordRepository interface {
	GetByID(ctx context.Context, id deliver.DeliverRecordID) (*deliver.DeliverRecord, error)
	Add(ctx context.Context, record *deliver.DeliverRecord) error
	Update(ctx context.Context, record *deliver.DeliverRecord) error
}
