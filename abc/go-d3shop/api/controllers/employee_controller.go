package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-d3shop/application/commands"
	"github.com/yourusername/go-d3shop/application/integration_events"
	"github.com/yourusername/go-d3shop/application/services"
	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/pkg/events"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// EmployeeController 员工控制器
type EmployeeController struct {
	mediator        *mediator.MediatorV2
	employeeService *services.EmployeeService
	eventPublisher  events.IIntegrationEventPublisher
}

// NewEmployeeController 创建员工控制器
func NewEmployeeController(
	mediator *mediator.MediatorV2,
	employeeService *services.EmployeeService,
	eventPublisher events.IIntegrationEventPublisher,
) *EmployeeController {
	return &EmployeeController{
		mediator:        mediator,
		employeeService: employeeService,
		eventPublisher:  eventPublisher,
	}
}

// CreateEmployeeRequest 创建员工请求DTO
type CreateEmployeeRequestDTO struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DepartmentID int64  `json:"departmentId" binding:"required,min=1"`
	Position     string `json:"position" binding:"required"`
	BaseSalary   int    `json:"baseSalary" binding:"required,min=1"`
}

// CreateEmployee 创建员工 - 使用MediatorV2发送命令
// @Summary 创建新员工
// @Description 创建新员工并触发相关的领域事件和集成事件
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param employee body CreateEmployeeRequestDTO true "员工信息"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Router /employees [post]
func (c *EmployeeController) CreateEmployee(ctx *gin.Context) {
	var req CreateEmployeeRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 创建命令
	cmd := &commands.CreateEmployeeCommand{
		Name:         req.Name,
		Email:        req.Email,
		DepartmentID: req.DepartmentID,
		Position:     req.Position,
		BaseSalary:   req.BaseSalary,
	}

	// 通过MediatorV2发送命令
	result, err := c.mediator.Send(ctx.Request.Context(), "CreateEmployee", cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	employeeID, ok := result.(employee.EmployeeID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid result type",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"employeeId": employeeID.Value(),
		"message":    "员工创建成功",
	})
}

// CreateEmployeeUsingService 创建员工 - 使用服务层（展示另一种方式）
// @Summary 创建新员工（服务层方式）
// @Description 通过服务层创建新员工
// @Tags 员工管理
// @Accept json
// @Produce json
// @Param employee body CreateEmployeeRequestDTO true "员工信息"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Router /employees/service [post]
func (c *EmployeeController) CreateEmployeeUsingService(ctx *gin.Context) {
	var req CreateEmployeeRequestDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 直接调用服务层
	employeeID, err := c.employeeService.CreateEmployee(ctx.Request.Context(), services.CreateEmployeeRequest{
		Name:         req.Name,
		Email:        req.Email,
		DepartmentID: req.DepartmentID,
		Position:     req.Position,
		BaseSalary:   req.BaseSalary,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"employeeId": employeeID.Value(),
		"message":    "员工创建成功（通过服务层）",
	})
}

// GetEmployee 获取员工信息
// @Summary 获取员工信息
// @Description 根据ID获取员工详细信息
// @Tags 员工管理
// @Produce json
// @Param id path int true "员工ID"
// @Success 200 {object} map[string]interface{} "员工信息"
// @Router /employees/{id} [get]
func (c *EmployeeController) GetEmployee(ctx *gin.Context) {
	employeeIDStr := ctx.Param("id")
	employeeIDInt, err := strconv.ParseInt(employeeIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的员工ID",
		})
		return
	}

	// 这里应该有查询命令或查询服务
	// 为了演示，我们暂时返回模拟数据
	ctx.JSON(http.StatusOK, gin.H{
		"employeeId": employeeIDInt,
		"name":       "示例员工",
		"email":      "example@company.com",
		"department": "技术部",
		"position":   "工程师",
	})
}

// SimulateEmployeeJoinedEvent 模拟员工入职集成事件（用于测试）
// @Summary 模拟员工入职事件
// @Description 模拟发送员工入职集成事件到消息队列
// @Tags 员工管理
// @Produce json
// @Param id path int true "员工ID"
// @Success 200 {object} map[string]interface{} "发送成功"
// @Router /employees/{id}/simulate-event [post]
func (c *EmployeeController) SimulateEmployeeJoinedEvent(ctx *gin.Context) {
	employeeIDStr := ctx.Param("id")
	employeeIDInt, err := strconv.ParseInt(employeeIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的员工ID",
		})
		return
	}

	// 创建模拟的集成事件
	event := integration_events.NewEmployeeJoinedIntegrationEvent(
		employee.NewEmployeeID(employeeIDInt),
		"测试员工",
		"test@company.com",
		employee.NewDepartmentID(1),
		"技术部",
		employee.NewEmployeeID(999),
		"manager@company.com",
		"高级工程师",
	)

	// 发布事件到消息队列
	err = c.eventPublisher.PublishAsync(ctx.Request.Context(), event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "员工入职事件已发送",
	})
}

// RegisterRoutes 注册路由
func (c *EmployeeController) RegisterRoutes(router *gin.RouterGroup) {
	employeeGroup := router.Group("/employees")
	{
		employeeGroup.POST("", c.CreateEmployee)                                 // 使用MediatorV2
		employeeGroup.POST("/service", c.CreateEmployeeUsingService)             // 使用服务层
		employeeGroup.GET("/:id", c.GetEmployee)                                 // 获取员工信息
		employeeGroup.POST("/:id/simulate-event", c.SimulateEmployeeJoinedEvent) // 模拟事件
	}
}
