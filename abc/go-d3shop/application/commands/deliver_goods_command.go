package commands

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/deliver"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// DeliverGoodsCommand 发货命令
type DeliverGoodsCommand struct {
	OrderID order.OrderID
}

// 确保实现ICommand接口
var _ mediator.ICommand = (*DeliverGoodsCommand)(nil)

// DeliverGoodsCommandHandler 发货命令处理器
type DeliverGoodsCommandHandler struct {
	deliverRepo repositories.IDeliverRecordRepository
}

// NewDeliverGoodsCommandHandler 创建命令处理器
func NewDeliverGoodsCommandHandler(deliverRepo repositories.IDeliverRecordRepository) *DeliverGoodsCommandHandler {
	return &DeliverGoodsCommandHandler{
		deliverRepo: deliverRepo,
	}
}

// Handle 处理命令
func (h *DeliverGoodsCommandHandler) Handle(ctx context.Context, cmd DeliverGoodsCommand) (deliver.DeliverRecordID, error) {
	// 创建发货记录
	record := deliver.NewDeliverRecord(cmd.OrderID)

	// 保存到仓储
	err := h.deliverRepo.Add(ctx, record)
	if err != nil {
		return deliver.DeliverRecordID{}, err
	}

	return record.ID, nil
}
