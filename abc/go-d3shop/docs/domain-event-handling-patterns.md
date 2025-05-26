# 领域事件处理模式

## 概述

在MediatorV2中，有两种处理领域事件的方式：

1. **事件转命令模式**（Event-to-Command）：保持与原MediatR一致的风格
2. **直接处理模式**（Direct Handling）：更符合Go的简洁风格

## 方式一：事件转命令模式

### 流程图
```
订单创建 → OrderCreatedDomainEvent → EventHandler → DeliverGoodsCommand → CommandHandler → 创建发货记录
```

### 实现示例
```go
// 1. 事件处理器将事件转换为命令
med.RegisterNotificationHandler("OrderCreated", func(ctx context.Context, notification interface{}) error {
    event := notification.(*events.OrderCreatedDomainEvent)
    orderAgg := event.Order.(*order.Order)
    
    // 转换为命令
    cmd := &commands.DeliverGoodsCommand{
        OrderID: orderAgg.ID,
    }
    
    // 发送命令
    _, err := med.Send(ctx, "DeliverGoods", cmd)
    return err
})

// 2. 命令处理器执行业务逻辑
med.RegisterCommandHandler("DeliverGoods", func(ctx context.Context, request interface{}) (interface{}, error) {
    cmd := request.(*commands.DeliverGoodsCommand)
    
    record := deliver.NewDeliverRecord(cmd.OrderID)
    err := deliverRepo.Add(ctx, record)
    
    return record.ID, err
})
```

### 优点
- **解耦性强**：事件和命令分离，各司其职
- **可测试性好**：可以单独测试命令处理器
- **复用性高**：命令可以被多个地方调用
- **符合CQRS**：保持命令查询职责分离

### 缺点
- **代码量多**：需要定义额外的命令和处理器
- **性能开销**：多一次中介者调用
- **调试复杂**：调用链更长

### 适用场景
- 复杂的业务流程
- 需要事务编排的场景
- 命令需要被多处复用
- 团队熟悉CQRS模式

## 方式二：直接处理模式

### 流程图
```
订单创建 → OrderCreatedDomainEvent → EventHandler → 直接创建发货记录
```

### 实现示例
```go
// 直接在事件处理器中执行业务逻辑
med.RegisterNotificationHandler("OrderCreated", func(ctx context.Context, notification interface{}) error {
    event := notification.(*events.OrderCreatedDomainEvent)
    orderAgg := event.Order.(*order.Order)
    
    // 直接执行业务逻辑
    record := deliver.NewDeliverRecord(orderAgg.ID)
    return deliverRepo.Add(ctx, record)
})
```

### 优点
- **简洁直接**：代码量少，易于理解
- **性能更好**：减少中介者调用
- **调试简单**：调用链短
- **符合Go风格**：简单直接

### 缺点
- **耦合度高**：事件处理器直接依赖仓储
- **复用性差**：逻辑无法被其他地方调用
- **测试复杂**：需要mock更多依赖

### 适用场景
- 简单的业务逻辑
- 一对一的事件响应
- 性能敏感的场景
- 团队偏好Go风格

## 混合使用建议

在实际项目中，可以根据具体场景混合使用两种模式：

```go
func RegisterEventHandlers(med *mediator.MediatorV2) {
    // 简单逻辑：直接处理
    med.RegisterNotificationHandler("UserRegistered", func(ctx context.Context, event interface{}) error {
        // 直接创建用户配置
        return userConfigRepo.CreateDefault(event.(*UserRegisteredEvent).UserID)
    })
    
    // 复杂流程：转换为命令
    med.RegisterNotificationHandler("OrderPaid", func(ctx context.Context, event interface{}) error {
        e := event.(*OrderPaidEvent)
        
        // 触发多个命令
        commands := []struct{
            name string
            cmd  interface{}
        }{
            {"UpdateInventory", &UpdateInventoryCommand{OrderID: e.OrderID}},
            {"SendNotification", &SendNotificationCommand{UserID: e.UserID}},
            {"UpdateStatistics", &UpdateStatisticsCommand{Amount: e.Amount}},
        }
        
        for _, c := range commands {
            if _, err := med.Send(ctx, c.name, c.cmd); err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

## 决策矩阵

| 因素 | 事件转命令 | 直接处理 |
|------|-----------|----------|
| 业务复杂度 | 高 ✓ | 低 ✓ |
| 代码复用需求 | 高 ✓ | 低 ✓ |
| 性能要求 | 一般 | 高 ✓ |
| 团队背景 | DDD/CQRS ✓ | Go ✓ |
| 测试策略 | 单元测试 ✓ | 集成测试 ✓ |
| 事务需求 | 复杂 ✓ | 简单 ✓ |

## 最佳实践

1. **保持一致性**：在同一个限界上下文内保持一致的处理风格
2. **文档化决策**：记录为什么选择某种模式
3. **定期重构**：随着业务发展，适时调整处理模式
4. **性能监控**：监控事件处理的性能，必要时优化

## 示例：订单系统的事件处理

```go
// main.go中的配置
useEventToCommandStyle := os.Getenv("EVENT_HANDLING_STYLE") == "command"

if useEventToCommandStyle {
    // 使用事件转命令风格
    domain_event_handlers.RegisterDomainEventHandlers(mediatorV2, repos...)
} else {
    // 使用直接处理风格
    orderService.RegisterEventHandlers()
}
```

这样可以通过环境变量灵活切换处理模式，便于对比和迁移。 