# 员工入职流程 - 领域事件与集成事件示例

## 业务场景

当新员工入职时，系统需要：
1. 创建员工记录（Employee上下文）
2. 更新部门员工数（Department上下文）- **领域事件**
3. 创建默认工资记录（Salary上下文）- **领域事件**
4. 通知部门经理（Notification服务）- **集成事件**

## 事件流程图

```
员工入职
    │
    ├─> 保存员工记录
    │
    ├─> 发布 EmployeeJoinedDomainEvent（领域事件）
    │   │
    │   ├─> 处理器1：更新部门员工数（同步，同事务）
    │   │
    │   └─> 处理器2：创建工资记录（同步，同事务）
    │
    └─> 发布 EmployeeJoinedIntegrationEvent（集成事件）
        │
        └─> 通知服务：发送邮件给部门经理（异步，跨服务）
```

## 代码实现

### 1. 员工聚合根触发领域事件

```go
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

### 2. 服务层处理流程

```go
func (s *EmployeeService) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (employee.EmployeeID, error) {
    // 1. 创建员工
    emp, err := employee.NewEmployee(...)
    
    // 2. 保存员工（触发领域事件）
    s.employeeRepo.Add(ctx, emp)
    
    // 3. 发布领域事件（同步处理）
    for _, event := range emp.GetDomainEvents() {
        s.mediator.Publish(ctx, event.EventName(), event)
    }
    
    // 4. 发布集成事件（异步处理）
    integrationEvent := integration_events.NewEmployeeJoinedIntegrationEvent(...)
    s.eventPublisher.PublishAsync(ctx, integrationEvent)
    
    return emp.ID, nil
}
```

### 3. 领域事件处理器

```go
// 直接处理风格
s.mediator.RegisterNotificationHandler("EmployeeJoined", func(ctx context.Context, notification interface{}) error {
    event := notification.(*events.EmployeeJoinedDomainEvent)
    emp := event.Employee.(*employee.Employee)
    
    // 更新部门员工数
    dept, _ := s.departmentRepo.GetByID(ctx, emp.DepartmentID)
    dept.AddEmployee()
    s.departmentRepo.Update(ctx, dept)
    
    // 创建工资记录
    salaryRecord := salary.NewSalaryRecord(emp.ID, emp.BaseSalary)
    s.salaryRepo.Add(ctx, salaryRecord)
    
    return nil
})
```

### 4. 集成事件处理器

```go
func (h *EmployeeJoinedIntegrationEventHandler) HandleAsync(ctx context.Context, event *EmployeeJoinedIntegrationEvent) error {
    // 发送邮件给部门经理
    if event.ManagerEmail != "" {
        h.sendEmailToManager(event)
    }
    
    // 发送欢迎邮件给新员工
    h.sendWelcomeEmail(event)
    
    return nil
}
```

## 关键区别

### 领域事件
- **范围**：限于本地限界上下文
- **处理**：同步，在同一事务内
- **数据**：包含完整的领域对象
- **失败**：导致整个操作回滚
- **示例**：更新部门员工数、创建工资记录

### 集成事件
- **范围**：跨服务/限界上下文
- **处理**：异步，通过消息队列
- **数据**：只包含必要的数据
- **失败**：不影响主流程，可重试
- **示例**：发送通知邮件

## 测试运行

运行测试查看完整流程：

```bash
cd go-d3shop
go test -v ./test -run TestEmployeeJoinedFlow
```

输出示例：
```
=== 开始测试员工入职流程 ===
当前部门员工数: 5

1. 创建新员工...
员工创建成功，ID: 1
工资记录已创建: ID=1, 员工ID=1, 基本工资=15000
集成事件已发布: EmployeeJoinedIntegrationEvent

2. 验证领域事件处理结果...
部门员工数已更新: 5 -> 6
工资记录数: 1
工资记录详情: 基本工资=15000, 月份=2024-01

3. 集成事件发布数: 1

=== 通知服务：处理员工入职事件 ===
新员工信息：
  姓名: 李四
  邮箱: lisi@company.com
  部门: 技术部
  职位: 高级工程师
发送邮件到: manager@company.com
发送欢迎邮件到: lisi@company.com

=== 员工入职流程测试完成 ===
```

## 设计考虑

1. **事务边界**：领域事件在同一事务内处理，保证数据一致性
2. **性能优化**：集成事件异步处理，不阻塞主流程
3. **错误处理**：领域事件失败会回滚，集成事件失败可重试
4. **解耦设计**：通过事件实现限界上下文之间的松耦合 