package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-d3shop/api/controllers"
	"github.com/yourusername/go-d3shop/application/commands"
	"github.com/yourusername/go-d3shop/application/domain_event_handlers"
	"github.com/yourusername/go-d3shop/application/services"
	"github.com/yourusername/go-d3shop/domain/aggregates/department"
	"github.com/yourusername/go-d3shop/domain/aggregates/employee"
	"github.com/yourusername/go-d3shop/infrastructure/messaging"
	"github.com/yourusername/go-d3shop/infrastructure/persistence"
	"github.com/yourusername/go-d3shop/infrastructure/repositories"
	"github.com/yourusername/go-d3shop/pkg/mediator"
	"github.com/yourusername/go-d3shop/pkg/mediator/behaviors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// LoggingBehavior 日志管道行为
type LoggingBehavior struct{}

// Handle 处理请求
func (b *LoggingBehavior) Handle(ctx context.Context, request mediator.IRequest, next mediator.RequestHandlerFunc) (interface{}, error) {
	log.Printf("开始处理请求: %T", request)

	result, err := next(ctx)

	if err != nil {
		log.Printf("请求处理失败: %v", err)
	} else {
		log.Printf("请求处理成功")
	}

	return result, err
}

// App 应用程序结构体 - 使用MediatorV2
type App struct {
	DB              *gorm.DB
	MediatorV2      *mediator.MediatorV2
	OrderService    *services.OrderServiceV2
	EmployeeService *services.EmployeeService
	EventPublisher  *messaging.RabbitMQPublisher
}

// NewApp 创建应用程序实例
func NewApp() (*App, error) {
	// 初始化数据库
	dsn := "user:password@tcp(127.0.0.1:3306)/d3shop?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	// db.AutoMigrate(&order.Order{}, &deliver.DeliverRecord{}, &employee.Employee{}, &department.Department{}, &salary.SalaryRecord{})

	// 创建MediatorV2
	mediatorV2 := mediator.NewMediatorV2()

	// 添加管道行为
	mediatorV2.AddPipelineBehavior(&LoggingBehavior{})
	mediatorV2.AddPipelineBehavior(behaviors.NewValidationBehavior())

	// 创建旧的mediator用于兼容DbContext
	oldMediator := mediator.NewMediator()

	// 创建数据库上下文
	dbContext := persistence.NewDbContext(db, oldMediator)

	// 创建仓储
	orderRepo := repositories.NewOrderRepository(dbContext)
	deliverRepo := repositories.NewDeliverRecordRepository(dbContext)
	employeeRepo := repositories.NewEmployeeRepository(dbContext)
	departmentRepo := repositories.NewDepartmentRepository(dbContext)
	salaryRepo := repositories.NewSalaryRepository(dbContext)

	// 注册命令处理器
	commands.RegisterCommandHandlers(mediatorV2, orderRepo)

	// 注册发货命令处理器（用于事件转命令模式）
	domain_event_handlers.RegisterDeliverGoodsCommandHandler(mediatorV2, deliverRepo)

	// 注册员工相关命令处理器
	commands.RegisterEmployeeCommandHandlers(mediatorV2, employeeRepo, departmentRepo, salaryRepo)

	// 创建应用服务
	orderService := services.NewOrderServiceV2(orderRepo, deliverRepo, mediatorV2)

	// 创建集成事件发布器
	eventPublisher, err := messaging.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/", "d3shop-events")
	if err != nil {
		log.Printf("Warning: Failed to create RabbitMQ publisher: %v", err)
		// 使用模拟的发布器
		eventPublisher = nil
	}

	// 创建员工服务
	employeeService := services.NewEmployeeService(
		employeeRepo,
		departmentRepo,
		salaryRepo,
		mediatorV2,
		eventPublisher,
	)

	// 注册领域事件处理器
	useCommandStyle := false // 可以通过配置切换

	// 订单事件处理器
	orderService.RegisterEventHandlers()

	// 员工事件处理器
	domain_event_handlers.RegisterEmployeeDomainEventHandlers(
		mediatorV2,
		departmentRepo,
		salaryRepo,
		useCommandStyle,
	)

	// 初始化测试数据
	initTestData(db)

	return &App{
		DB:              db,
		MediatorV2:      mediatorV2,
		OrderService:    orderService,
		EmployeeService: employeeService,
		EventPublisher:  eventPublisher,
	}, nil
}

// initTestData 初始化测试数据
func initTestData(db *gorm.DB) {
	// 创建测试部门
	dept := &department.Department{
		ID:            employee.NewDepartmentID(1),
		Name:          "技术部",
		ManagerID:     employee.NewEmployeeID(999),
		EmployeeCount: 10,
	}
	db.FirstOrCreate(dept, "id = ?", 1)

	// 创建部门经理
	manager := &employee.Employee{
		ID:           employee.NewEmployeeID(999),
		Name:         "张经理",
		Email:        "manager@company.com",
		DepartmentID: employee.NewDepartmentID(1),
		Position:     "技术总监",
		BaseSalary:   30000,
	}
	db.FirstOrCreate(manager, "id = ?", 999)

	log.Println("测试数据初始化完成")
}

// Run 运行应用程序
func (app *App) Run() error {
	// 创建Gin路由
	r := gin.Default()

	// 创建控制器
	orderController := controllers.NewOrderControllerV3(app.MediatorV2, app.OrderService, app.EventPublisher)
	employeeController := controllers.NewEmployeeController(app.MediatorV2, app.EmployeeService, app.EventPublisher)

	// 注册路由
	api := r.Group("/api")
	orderController.RegisterRoutes(api)
	employeeController.RegisterRoutes(api)

	// 启动集成事件消费者
	go app.startEventConsumer()

	// 启动HTTP服务器
	log.Println("Server starting on :8080...")
	log.Println("使用MediatorV2版本 - 包含订单和员工模块")
	log.Println("")
	log.Println("API路径:")
	log.Println("订单模块:")
	log.Println("  POST   /api/orders              - 创建订单")
	log.Println("  POST   /api/orders/service      - 创建订单（服务层）")
	log.Println("  PUT    /api/orders/:id/pay      - 支付订单")
	log.Println("  POST   /api/orders/:id/pay-event - 发送支付事件")
	log.Println("")
	log.Println("员工模块:")
	log.Println("  POST   /api/employees            - 创建员工（使用MediatorV2）")
	log.Println("  POST   /api/employees/service    - 创建员工（使用服务层）")
	log.Println("  GET    /api/employees/:id        - 获取员工信息")
	log.Println("  POST   /api/employees/:id/simulate-event - 模拟员工入职事件")

	return r.Run(":8080")
}

// startEventConsumer 启动事件消费者
func (app *App) startEventConsumer() {
	log.Println("Event consumer started")

	// 在实际应用中，这里应该连接到RabbitMQ并消费消息
	// 这里只是示例，展示如何处理不同类型的集成事件
}

// Cleanup 清理资源
func (app *App) Cleanup() {
	if app.EventPublisher != nil {
		app.EventPublisher.Close()
	}
}

func main() {
	// 创建应用
	app, err := NewApp()
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}
	defer app.Cleanup()

	// 运行应用
	if err := app.Run(); err != nil {
		log.Fatal("Failed to run app:", err)
	}
}
