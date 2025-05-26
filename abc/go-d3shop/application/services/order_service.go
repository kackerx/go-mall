package services

import (
	"context"
	"errors"

	"github.com/yourusername/go-d3shop/domain/aggregates/deliver"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/domain/events"
	"github.com/yourusername/go-d3shop/domain/repositories"
	"github.com/yourusername/go-d3shop/pkg/cqrs"
)

// OrderService 订单服务 - Go风格的应用服务
type OrderService struct {
	orderRepo   repositories.IOrderRepository
	deliverRepo repositories.IDeliverRecordRepository
	eventBus    *cqrs.EventBus
}

// NewOrderService 创建订单服务
func NewOrderService(
	orderRepo repositories.IOrderRepository,
	deliverRepo repositories.IDeliverRecordRepository,
	eventBus *cqrs.EventBus,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		deliverRepo: deliverRepo,
		eventBus:    eventBus,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Name  string
	Price int
	Count int
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (order.OrderID, error) {
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

	// 处理领域事件 - 直接在服务中处理，而不是通过复杂的事件总线
	// 这是Go中更直接的方式
	for _, event := range orderAgg.GetDomainEvents() {
		switch e := event.(type) {
		case *events.OrderCreatedDomainEvent:
			// 直接调用发货逻辑
			if err := s.createDeliveryRecord(ctx, orderAgg.ID); err != nil {
				// 记录错误但不影响主流程
				// 在实际应用中应该有更好的错误处理
			}
		default:
			// 忽略未知事件
			_ = e
		}
	}
	orderAgg.ClearDomainEvents()

	return orderAgg.ID, nil
}

// PayOrder 支付订单
func (s *OrderService) PayOrder(ctx context.Context, orderID order.OrderID) error {
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
	return s.orderRepo.Update(ctx, orderAgg)
}

// createDeliveryRecord 创建发货记录（内部方法）
func (s *OrderService) createDeliveryRecord(ctx context.Context, orderID order.OrderID) error {
	record := deliver.NewDeliverRecord(orderID)
	return s.deliverRepo.Add(ctx, record)
}
