package commands

import (
	"context"

	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/cqrs"
)

// CreateOrderCommandV2 创建订单命令（Go风格版本）
type CreateOrderCommandV2 struct {
	Name  string
	Price int
	Count int
}

// CommandName 实现Command接口
func (c CreateOrderCommandV2) CommandName() string {
	return "CreateOrder"
}

// 确保实现Command接口
var _ cqrs.Command = (*CreateOrderCommandV2)(nil)

// CreateOrderHandlerV2 创建订单处理器
type CreateOrderHandlerV2 struct {
	orderRepo repositories.IOrderRepository
	eventBus  *cqrs.EventBus
}

// NewCreateOrderHandlerV2 创建处理器
func NewCreateOrderHandlerV2(orderRepo repositories.IOrderRepository, eventBus *cqrs.EventBus) *CreateOrderHandlerV2 {
	return &CreateOrderHandlerV2{
		orderRepo: orderRepo,
		eventBus:  eventBus,
	}
}

// Handle 处理命令
func (h *CreateOrderHandlerV2) Handle(ctx context.Context, cmd cqrs.Command) (interface{}, error) {
	// 类型断言
	createCmd, ok := cmd.(*CreateOrderCommandV2)
	if !ok {
		return nil, cqrs.ErrInvalidCommand
	}

	// 创建订单聚合根
	orderAgg := order.NewOrder(createCmd.Name, createCmd.Count)

	// 保存到仓储
	err := h.orderRepo.Add(ctx, orderAgg)
	if err != nil {
		return nil, err
	}

	// 发布领域事件
	for _, event := range orderAgg.GetDomainEvents() {
		if domainEvent, ok := event.(cqrs.Event); ok {
			if err := h.eventBus.Publish(ctx, domainEvent); err != nil {
				// 记录错误但不影响主流程
				continue
			}
		}
	}
	orderAgg.ClearDomainEvents()

	return orderAgg.ID, nil
}
