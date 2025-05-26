package services

import (
	"context"
	"errors"

	"github.com/yourusername/go-d3shop/domain/aggregates/deliver"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/events"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// OrderServiceV2 订单服务 - 使用MediatorV2
type OrderServiceV2 struct {
	orderRepo   repositories.IOrderRepository
	deliverRepo repositories.IDeliverRecordRepository
	mediator    *mediator.MediatorV2
}

// NewOrderServiceV2 创建订单服务
func NewOrderServiceV2(
	orderRepo repositories.IOrderRepository,
	deliverRepo repositories.IDeliverRecordRepository,
	mediator *mediator.MediatorV2,
) *OrderServiceV2 {
	return &OrderServiceV2{
		orderRepo:   orderRepo,
		deliverRepo: deliverRepo,
		mediator:    mediator,
	}
}

// CreateOrder 创建订单
func (s *OrderServiceV2) CreateOrder(ctx context.Context, req CreateOrderRequest) (order.OrderID, error) {
	// 验证请求
	if req.Name == "" {
		return order.OrderID{}, errors.New("订单名称不能为空")
	}
	if req.Count <= 0 {
		return order.OrderID{}, errors.New("订单数量必须大于0")
	}

	// 创建订单聚合根
	orderAgg := order.NewOrder(req.Name, req.Count)

	// 保存订单
	if err := s.orderRepo.Add(ctx, orderAgg); err != nil {
		return order.OrderID{}, err
	}

	// 通过MediatorV2发布领域事件
	for _, event := range orderAgg.GetDomainEvents() {
		// 获取事件名称并发布
		if domainEvent, ok := event.(interface{ EventName() string }); ok {
			err := s.mediator.Publish(ctx, domainEvent.EventName(), event)
			if err != nil {
				// 记录错误但不影响主流程
				// 在实际应用中应该有更好的错误处理
				continue
			}
		}
	}
	orderAgg.ClearDomainEvents()

	return orderAgg.ID, nil
}

// PayOrder 支付订单
func (s *OrderServiceV2) PayOrder(ctx context.Context, orderID order.OrderID) error {
	// 获取订单
	orderAgg, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if orderAgg == nil {
		return errors.New("订单不存在")
	}

	// 执行支付
	if err := orderAgg.OrderPaid(); err != nil {
		return err
	}

	// 更新订单
	if err := s.orderRepo.Update(ctx, orderAgg); err != nil {
		return err
	}

	// 发布领域事件
	for _, event := range orderAgg.GetDomainEvents() {
		if domainEvent, ok := event.(interface{ EventName() string }); ok {
			err := s.mediator.Publish(ctx, domainEvent.EventName(), event)
			if err != nil {
				continue
			}
		}
	}
	orderAgg.ClearDomainEvents()

	return nil
}

// RegisterEventHandlers 注册事件处理器
func (s *OrderServiceV2) RegisterEventHandlers() {
	// 注册订单创建事件处理器
	s.mediator.RegisterNotificationHandler("OrderCreated", func(ctx context.Context, notification interface{}) error {
		if event, ok := notification.(*events.OrderCreatedDomainEvent); ok {
			// 从事件中获取订单信息
			if orderAgg, ok := event.Order.(*order.Order); ok {
				// 创建发货记录
				record := deliver.NewDeliverRecord(orderAgg.ID)
				return s.deliverRepo.Add(ctx, record)
			}
		}
		return nil
	})

	// 注册订单支付事件处理器
	s.mediator.RegisterNotificationHandler("OrderPaid", func(ctx context.Context, notification interface{}) error {
		// 这里可以添加订单支付后的其他处理逻辑
		// 比如发送通知、更新统计等
		return nil
	})
}
