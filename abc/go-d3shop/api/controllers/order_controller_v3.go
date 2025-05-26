package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-d3shop/application/commands"
	"github.com/yourusername/go-d3shop/application/integration_events"
	"github.com/yourusername/go-d3shop/application/services"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/pkg/events"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// OrderControllerV3 订单控制器（使用MediatorV2）
type OrderControllerV3 struct {
	mediator       *mediator.MediatorV2
	orderService   *services.OrderServiceV2
	eventPublisher events.IIntegrationEventPublisher
}

// NewOrderControllerV3 创建订单控制器
func NewOrderControllerV3(
	mediator *mediator.MediatorV2,
	orderService *services.OrderServiceV2,
	eventPublisher events.IIntegrationEventPublisher,
) *OrderControllerV3 {
	return &OrderControllerV3{
		mediator:       mediator,
		orderService:   orderService,
		eventPublisher: eventPublisher,
	}
}

// CreateOrderRequestDTO 创建订单请求DTO
type CreateOrderRequestDTOV3 struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required,min=1"`
	Count int    `json:"count" binding:"required,min=1"`
}

// CreateOrder 创建订单 - 使用MediatorV2发送命令
func (c *OrderControllerV3) CreateOrder(ctx *gin.Context) {
	var req CreateOrderRequestDTOV3
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 创建命令
	cmd := &commands.CreateOrderCommandV3{
		Name:  req.Name,
		Price: req.Price,
		Count: req.Count,
	}

	// 通过MediatorV2发送命令
	result, err := c.mediator.Send(ctx.Request.Context(), "CreateOrder", cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	orderID, ok := result.(order.OrderID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid result type",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orderId": orderID.Value(),
		"message": "订单创建成功",
	})
}

// CreateOrderUsingService 创建订单 - 使用服务层（展示另一种方式）
func (c *OrderControllerV3) CreateOrderUsingService(ctx *gin.Context) {
	var req CreateOrderRequestDTOV3
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 直接调用服务层
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
		"message": "订单创建成功（通过服务层）",
	})
}

// PayOrder 支付订单 - 使用MediatorV2
func (c *OrderControllerV3) PayOrder(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderIDInt, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的订单ID",
		})
		return
	}

	// 创建命令
	cmd := &commands.OrderPaidCommandV2{
		OrderID: order.NewOrderID(orderIDInt),
	}

	// 通过MediatorV2发送命令
	_, err = c.mediator.Send(ctx.Request.Context(), "PayOrder", cmd)
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
func (c *OrderControllerV3) SendPaymentEvent(ctx *gin.Context) {
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
func (c *OrderControllerV3) RegisterRoutes(router *gin.RouterGroup) {
	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("", c.CreateOrder)                     // 使用MediatorV2
		orderGroup.POST("/service", c.CreateOrderUsingService) // 使用服务层
		orderGroup.PUT("/:id/pay", c.PayOrder)
		orderGroup.POST("/:id/pay-event", c.SendPaymentEvent)
	}
}
