# 员工入职完整流程 - 从API到事件处理

## 架构概览

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐     ┌──────────────┐
│   客户端    │────▶│  Controller  │────▶│  MediatorV2 │────▶│   Command    │
└─────────────┘     └──────────────┘     └─────────────┘     │   Handler    │
                                                              └──────────────┘
                                                                      │
                                                                      ▼
┌─────────────┐     ┌──────────────┐     ┌─────────────┐     ┌──────────────┐
│  集成事件   │◀────│   Service    │◀────│ Repository  │◀────│   Domain     │
│  发布器     │     │              │     │             │     │   Model      │
└─────────────┘     └──────────────┘     └─────────────┘     └──────────────┘
                            │
                            ▼
                    ┌──────────────┐
                    │ Domain Event │
                    │  Handlers    │
                    └──────────────┘
```

## 完整调用流程

### 1. API入口 (Controller层)

```go
// POST /api/employees
func (c *EmployeeController) CreateEmployee(ctx *gin.Context) {
    // 1. 接收请求
    var req CreateEmployeeRequestDTO
    ctx.ShouldBindJSON(&req)
    
    // 2. 创建命令
    cmd := &commands.CreateEmployeeCommand{
        Name:         req.Name,
        Email:        req.Email,
        DepartmentID: req.DepartmentID,
        Position:     req.Position,
        BaseSalary:   req.BaseSalary,
    }
    
    // 3. 通过MediatorV2发送命令
    result, err := c.mediator.Send(ctx, "CreateEmployee", cmd)
}
```

### 2. 命令处理 (Application层)

```go
// CreateEmployeeCommandHandler
func (h *CreateEmployeeCommandHandler) Handle(ctx context.Context, cmd CreateEmployeeCommand) (employee.EmployeeID, error) {
    // 1. 验证部门存在
    dept, err := h.departmentRepo.GetByID(ctx, employee.NewDepartmentID(cmd.DepartmentID))
    
    // 2. 创建员工聚合根（触发领域事件）
    emp, err := employee.NewEmployee(
        cmd.Name,
        cmd.Email,
        employee.NewDepartmentID(cmd.DepartmentID),
        cmd.Position,
        cmd.BaseSalary,
    )
    
    // 3. 保存到仓储
    err = h.employeeRepo.Add(ctx, emp)
}
```

### 3. 领域模型 (Domain层)

```go
// NewEmployee 创建新员工
func NewEmployee(name, email string, departmentID DepartmentID, position string, baseSalary int) (*Employee, error) {
    employee := &Employee{
        Name:         name,
        Email:        email,
        DepartmentID: departmentID,
        Position:     position,
        JoinDate:     time.Now(),
        BaseSalary:   baseSalary,
    }
    
    // 添加领域事件
    employee.AddDomainEvent(events.NewEmployeeJoinedDomainEvent(employee))
    
    return employee, nil
}
```

### 4. 仓储保存 (Infrastructure层)

```go
// EmployeeRepository.Add
func (r *EmployeeRepository) Add(ctx context.Context, emp *employee.Employee) error {
    // 1. 保存员工到数据库
    err := r.dbContext.DB().WithContext(ctx).Create(emp).Error
    
    // 2. 发布领域事件
    return r.dbContext.PublishDomainEvents(ctx, emp)
}
```

### 5. 领域事件处理 (Application层)

```go
// EmployeeJoinedDomainEventHandler
func (h *EmployeeJoinedDomainEventHandler) Handle(ctx context.Context, event *events.EmployeeJoinedDomainEvent) error {
    emp := event.Employee.(*employee.Employee)
    
    // 1. 更新部门员工数
    dept, _ := h.departmentRepo.GetByID(ctx, emp.DepartmentID)
    dept.AddEmployee()
    h.departmentRepo.Update(ctx, dept)
    
    // 2. 创建工资记录
    salaryRecord := salary.NewSalaryRecord(emp.ID, emp.BaseSalary)
    h.salaryRepo.Add(ctx, salaryRecord)
}
```

### 6. 服务层发布集成事件 (Application层)

```go
// EmployeeService.CreateEmployee
func (s *EmployeeService) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (employee.EmployeeID, error) {
    // ... 创建员工逻辑 ...
    
    // 发布集成事件
    integrationEvent := integration_events.NewEmployeeJoinedIntegrationEvent(
        emp.ID,
        emp.Name,
        emp.Email,
        dept.ID,
        dept.Name,
        dept.ManagerID,
        managerEmail,
        emp.Position,
    )
    
    s.eventPublisher.PublishAsync(ctx, integrationEvent)
}
```

### 7. 集成事件处理 (跨服务)

```go
// EmployeeJoinedIntegrationEventHandler
func (h *EmployeeJoinedIntegrationEventHandler) HandleAsync(ctx context.Context, event *EmployeeJoinedIntegrationEvent) error {
    // 发送邮件给部门经理
    h.sendEmailToManager(event)
    
    // 发送欢迎邮件给新员工
    h.sendWelcomeEmail(event)
}
```

## 两种使用方式对比

### 方式1：通过MediatorV2（推荐）

```go
// Controller
result, err := c.mediator.Send(ctx, "CreateEmployee", cmd)
```

**优点**：
- 解耦控制器和业务逻辑
- 支持管道行为（日志、验证、事务）
- 统一的命令处理模式

### 方式2：直接调用服务

```go
// Controller
employeeID, err := c.employeeService.CreateEmployee(ctx, req)
```

**优点**：
- 更直接，调试方便
- 减少抽象层次
- 适合简单场景

## 事件处理的两种风格

### 风格1：直接处理

```go
// 在事件处理器中直接执行业务逻辑
dept.AddEmployee()
h.departmentRepo.Update(ctx, dept)
```

### 风格2：事件转命令

```go
// 将事件转换为新的命令
h.mediator.Send(ctx, "UpdateDepartmentEmployeeCount", &UpdateDepartmentEmployeeCountCommand{
    DepartmentID: emp.DepartmentID,
    Delta: 1,
})
```

## 测试API

### 创建员工（使用MediatorV2）
```bash
curl -X POST http://localhost:8080/api/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "李四",
    "email": "lisi@company.com",
    "departmentId": 1,
    "position": "高级工程师",
    "baseSalary": 15000
  }'
```

### 创建员工（使用服务层）
```bash
curl -X POST http://localhost:8080/api/employees/service \
  -H "Content-Type: application/json" \
  -d '{
    "name": "王五",
    "email": "wangwu@company.com",
    "departmentId": 1,
    "position": "产品经理",
    "baseSalary": 18000
  }'
```

### 模拟集成事件
```bash
curl -X POST http://localhost:8080/api/employees/1/simulate-event
```

## 运行应用

```bash
# 确保MySQL和RabbitMQ已启动
# 修改main.go使用mainEmployee函数

cd go-d3shop
go run .
```

## 关键设计点

1. **命令处理器注册**：在应用启动时注册所有命令处理器
2. **领域事件发布**：在仓储层统一处理领域事件发布
3. **事务边界**：领域事件在同一事务内处理，集成事件异步处理
4. **错误处理**：领域事件失败会回滚事务，集成事件失败不影响主流程
5. **依赖注入**：通过构造函数注入所有依赖，便于测试和维护 