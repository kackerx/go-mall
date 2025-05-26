package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-d3shop/application/integration_events"
	"github.com/yourusername/go-d3shop/application/services"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/pkg/events"
)

// OrderControllerV2 订单控制器（Go风格版本）
type OrderControllerV2 struct {
	orderService   *services.OrderService
	eventPublisher events.IIntegrationEventPublisher
}

// NewOrderControllerV2 创建订单控制器
func NewOrderControllerV2(
	orderService *services.OrderService,
	eventPublisher events.IIntegrationEventPublisher,
) *OrderControllerV2 {
	return &OrderControllerV2{
		orderService:   orderService,
		eventPublisher: eventPublisher,
	}
}

// CreateOrderRequestDTO 创建订单请求DTO
type CreateOrderRequestDTO struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required,min=1"`
	Count int    `json:"count" binding:"required,min=1"`
}

// CreateOrder 创建订单
func (c *OrderControllerV2) CreateOrder(ctx *gin.Context) {
	var req CreateOrderRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 调用应用服务
	orderID, err := c.orderService.CreateOrder(ctx.Request.Context(), services.CreateOrderRequest{
		Name:  req.Name,
		Price: req.Price,
		Count: req.Count,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orderId": orderID.Value(),
		"message": "订单创建成功",
	})
}

// PayOrder 支付订单
func (c *OrderControllerV2) PayOrder(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderIDInt, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的订单ID",
		})
		return
	}

	// 调用应用服务
	err = c.orderService.PayOrder(ctx.Request.Context(), order.NewOrderID(orderIDInt))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "订单支付成功",
	})
}

// SendPaymentEvent 发送支付事件（模拟外部支付系统）
func (c *OrderControllerV2) SendPaymentEvent(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderIDInt, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的订单ID",
		})
		return
	}

	// 创建集成事件
	event := integration_events.NewOrderPaidIntegrationEvent(order.NewOrderID(orderIDInt))

	// 发布事件到消息队列
	err = c.eventPublisher.PublishAsync(ctx.Request.Context(), event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "支付事件已发送",
	})
}

// RegisterRoutes 注册路由
func (c *OrderControllerV2) RegisterRoutes(router *gin.RouterGroup) {
	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("", c.CreateOrder)
		orderGroup.PUT("/:id/pay", c.PayOrder)
		orderGroup.POST("/:id/pay-event", c.SendPaymentEvent)
	}
}
