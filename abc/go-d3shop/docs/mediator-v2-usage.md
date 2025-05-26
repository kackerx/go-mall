# MediatorV2 使用说明

## 概述

MediatorV2 是一个更符合Go风格的中介者实现，避免了大量反射，使用函数类型和字符串键来管理处理器。

## 核心组件

### 1. 命令/查询处理器 (handlers)

用于处理命令和查询，一对一映射。

```go
// 注册命令处理器
mediator.RegisterCommandHandler("CreateOrder", func(ctx context.Context, request interface{}) (interface{}, error) {
    cmd := request.(*CreateOrderCommand)
    // 处理逻辑
    return orderID, nil
})

// 发送命令
result, err := mediator.Send(ctx, "CreateOrder", createOrderCmd)
```

### 2. 通知处理器 (notificationHandlers)

用于处理事件/通知，支持多个处理器。

```go
// 注册多个事件处理器
mediator.RegisterNotificationHandler("OrderCreated", handler1)
mediator.RegisterNotificationHandler("OrderCreated", handler2)

// 发布事件
err := mediator.Publish(ctx, "OrderCreated", orderCreatedEvent)
```

### 3. 管道行为 (pipelineBehaviors)

用于添加横切关注点。

```go
// 添加管道行为
mediator.AddPipelineBehavior(&LoggingBehavior{})
mediator.AddPipelineBehavior(&ValidationBehavior{})
```

## 完整使用流程

### 1. 控制器层

```go
// POST /api/v3/orders
func (c *OrderControllerV3) CreateOrder(ctx *gin.Context) {
    // 1. 接收请求
    var req CreateOrderRequestDTO
    ctx.ShouldBindJSON(&req)
    
    // 2. 创建命令
    cmd := &CreateOrderCommandV3{
        Name:  req.Name,
        Price: req.Price,
        Count: req.Count,
    }
    
    // 3. 通过MediatorV2发送命令
    result, err := c.mediator.Send(ctx, "CreateOrder", cmd)
    
    // 4. 返回结果
    orderID := result.(order.OrderID)
    ctx.JSON(200, gin.H{"orderId": orderID.Value()})
}
```

### 2. 命令处理器

```go
// 在应用启动时注册
mediator.RegisterCommandHandler("CreateOrder", func(ctx context.Context, request interface{}) (interface{}, error) {
    cmd := request.(*CreateOrderCommandV3)
    
    // 创建订单聚合根
    orderAgg := order.NewOrder(cmd.Name, cmd.Count)
    
    // 保存到仓储（会触发领域事件）
    err := orderRepo.Add(ctx, orderAgg)
    
    return orderAgg.ID, err
})
```

### 3. 领域事件处理

```go
// OrderService中处理领域事件
func (s *OrderServiceV2) CreateOrder(ctx context.Context, req CreateOrderRequest) (order.OrderID, error) {
    // 创建订单
    orderAgg := order.NewOrder(req.Name, req.Count)
    
    // 保存订单
    s.orderRepo.Add(ctx, orderAgg)
    
    // 发布领域事件到MediatorV2
    for _, event := range orderAgg.GetDomainEvents() {
        if domainEvent, ok := event.(interface{ EventName() string }); ok {
            // 通过事件名称发布
            s.mediator.Publish(ctx, domainEvent.EventName(), event)
        }
    }
    
    return orderAgg.ID, nil
}
```

### 4. 事件处理器注册

```go
// 注册领域事件处理器
func (s *OrderServiceV2) RegisterEventHandlers() {
    // 订单创建事件 -> 创建发货记录
    s.mediator.RegisterNotificationHandler("OrderCreated", 
        func(ctx context.Context, notification interface{}) error {
            event := notification.(*events.OrderCreatedDomainEvent)
            orderAgg := event.Order.(*order.Order)
            
            // 创建发货记录
            record := deliver.NewDeliverRecord(orderAgg.ID)
            return s.deliverRepo.Add(ctx, record)
        })
    
    // 订单支付事件 -> 其他业务逻辑
    s.mediator.RegisterNotificationHandler("OrderPaid", 
        func(ctx context.Context, notification interface{}) error {
            // 处理订单支付后的逻辑
            return nil
        })
}
```

## 执行流程图

```
用户请求
    ↓
OrderController.CreateOrder()
    ↓
mediator.Send("CreateOrder", cmd)
    ↓
LoggingBehavior.Handle() ─────┐
    ↓                         │ 管道行为
ValidationBehavior.Handle() ──┘
    ↓
CreateOrderCommandHandler()
    ↓
orderRepo.Add() → 触发 OrderCreatedDomainEvent
    ↓
mediator.Publish("OrderCreated", event)
    ↓
OrderCreatedEventHandler1() ──┐
    ↓                         │ 并行执行
OrderCreatedEventHandler2() ──┘
```

## 优势

1. **类型安全**：虽然使用interface{}，但在处理器内部立即进行类型断言
2. **性能更好**：避免了反射的开销
3. **调试友好**：可以直接看到处理器函数，便于调试
4. **灵活性**：支持函数式编程风格

## 与原版MediatR的对比

| 特性 | MediatR (v1) | MediatorV2 |
|------|--------------|------------|
| 处理器注册 | 基于反射类型 | 基于字符串键 |
| 性能 | 较慢（反射） | 较快（直接调用） |
| 类型安全 | 编译时 | 运行时（需要断言） |
| 调试 | 困难 | 容易 |
| Go风格 | 否 | 是 |

## 最佳实践

1. **命名约定**：使用清晰的命令/事件名称
2. **错误处理**：在处理器中进行类型断言时要处理错误
3. **注册时机**：在应用启动时完成所有注册
4. **避免循环依赖**：合理组织代码结构 