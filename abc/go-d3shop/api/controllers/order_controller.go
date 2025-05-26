package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-d3shop/application/commands"
	"github.com/yourusername/go-d3shop/application/integration_events"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/pkg/events"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// OrderController 订单控制器
type OrderController struct {
	mediator       mediator.IMediator
	eventPublisher events.IIntegrationEventPublisher
}

// NewOrderController 创建订单控制器
func NewOrderController(mediator mediator.IMediator, eventPublisher events.IIntegrationEventPublisher) *OrderController {
	return &OrderController{
		mediator:       mediator,
		eventPublisher: eventPublisher,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	Name  string `json:"name" binding:"required"`
	Price int    `json:"price" binding:"required,min=1"`
	Count int    `json:"count" binding:"required,min=1"`
}

// CreateOrder 创建订单
func (c *OrderController) CreateOrder(ctx *gin.Context) {
	var req CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建命令
	cmd := commands.CreateOrderCommand{
		Name:  req.Name,
		Price: req.Price,
		Count: req.Count,
	}

	// 发送命令
	result, err := c.mediator.Send(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderID, ok := result.(order.OrderID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid result type"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orderId": orderID.Value(),
		"message": "订单创建成功",
	})
}

// SendPaymentEvent 发送支付事件（模拟）
func (c *OrderController) SendPaymentEvent(ctx *gin.Context) {
	orderIDStr := ctx.Param("id")
	orderIDInt, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	// 创建集成事件
	event := integration_events.NewOrderPaidIntegrationEvent(order.NewOrderID(orderIDInt))

	// 发布事件
	err = c.eventPublisher.PublishAsync(ctx.Request.Context(), event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "支付事件已发送",
	})
}

// RegisterRoutes 注册路由
func (c *OrderController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/orders", c.CreateOrder)
	router.POST("/orders/:id/pay-event", c.SendPaymentEvent)
}
