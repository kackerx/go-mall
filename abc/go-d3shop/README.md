# Go-D3Shop - DDD + 事件驱动架构示例

这是一个使用Go语言重构的D3Shop项目，展示了如何在Go中实现DDD（领域驱动设计）、事件驱动架构和MediatR模式。

## 架构概述

### 核心概念

1. **领域驱动设计 (DDD)**
   - 聚合根 (Aggregate Root)
   - 实体 (Entity)
   - 值对象 (Value Object)
   - 领域事件 (Domain Event)
   - 仓储 (Repository)

2. **事件驱动架构**
   - **领域事件**：在本地事务内同步处理，确保强一致性
   - **集成事件**：通过消息队列异步处理，实现最终一致性

3. **MediatR模式**
   - 命令 (Command) / 查询 (Query) 分离
   - 请求/响应模式
   - 管道行为 (Pipeline Behaviors)

## 项目结构

```
go-d3shop/
├── domain/                      # 领域层
│   ├── aggregates/             # 聚合根
│   │   ├── order/             # 订单聚合
│   │   └── deliver/           # 发货聚合
│   ├── events/                # 领域事件
│   └── repositories/          # 仓储接口
│
├── application/                # 应用层
│   ├── commands/              # 命令和处理器
│   ├── queries/               # 查询和处理器
│   ├── domain_event_handlers/ # 领域事件处理器
│   ├── integration_events/    # 集成事件
│   └── integration_event_handlers/ # 集成事件处理器
│
├── infrastructure/            # 基础设施层
│   ├── persistence/          # 数据持久化
│   ├── repositories/         # 仓储实现
│   └── messaging/            # 消息队列
│
├── api/                      # 表现层
│   └── controllers/          # HTTP控制器
│
└── pkg/                      # 共享包
    ├── ddd/                  # DDD基础设施
    ├── mediator/             # MediatR实现
    └── events/               # 事件基础设施
```

## 核心流程

### 1. 创建订单流程

```
用户请求 -> OrderController -> CreateOrderCommand -> CreateOrderCommandHandler
    -> Order聚合根 (触发OrderCreatedDomainEvent) -> OrderRepository保存
    -> 发布领域事件 -> OrderCreatedDomainEventHandler 
    -> DeliverGoodsCommand -> 创建发货记录
```

### 2. 订单支付流程（跨服务）

```
外部支付服务 -> RabbitMQ -> OrderPaidIntegrationEvent 
    -> OrderPaidIntegrationEventHandler -> OrderPaidCommand
    -> Order聚合根 (触发OrderPaidDomainEvent) -> 更新订单状态
```

## 关键设计模式

### 1. 强类型ID
```go
type OrderID struct {
    ddd.Int64StronglyTypedId
}
```
避免原始类型困扰，提供类型安全。

### 2. 聚合根基类
```go
type Order struct {
    ddd.BaseEntity  // 包含领域事件管理
    ID    OrderID
    // ... 其他属性
}
```

### 3. MediatR模式
```go
// 发送命令
result, err := mediator.Send(ctx, command)

// 发布事件
err := mediator.Publish(ctx, event)
```

### 4. 管道行为
- **ValidationBehavior**: 自动验证请求
- **UnitOfWorkBehavior**: 事务管理

## 事件处理机制

### 领域事件
- 在聚合根中触发：`order.AddDomainEvent(event)`
- 在仓储保存后同步发布
- 在同一事务内处理，保证强一致性

### 集成事件
- 通过RabbitMQ发布/订阅
- 跨服务边界通信
- 最终一致性保证

## 运行项目

1. 安装依赖：
```bash
go mod download
```

2. 启动MySQL和RabbitMQ：
```bash
docker-compose up -d mysql rabbitmq
```

3. 运行应用：
```bash
go run main.go
```

## API示例

### 创建订单
```bash
POST /api/orders
{
    "name": "商品名称",
    "price": 100,
    "count": 2
}
```

### 发送支付事件（模拟）
```bash
POST /api/orders/{orderId}/pay-event
```

## 与C#版本的对比

| C# 特性 | Go 实现 |
|---------|---------|
| MediatR | 自定义mediator包，使用反射实现 |
| Entity Framework | GORM |
| 依赖注入 | 手动依赖注入 |
| 接口 | interface{} |
| 泛型 | Go 1.18+ 泛型 |
| async/await | goroutine + context |
| CAP | RabbitMQ + 自定义实现 |

## 扩展建议

1. **添加更多聚合根**：如用户、商品等
2. **实现Saga模式**：处理分布式事务
3. **添加事件溯源**：Event Sourcing
4. **完善错误处理**：自定义错误类型
5. **添加日志和监控**：使用zap + prometheus
6. **实现CQRS读模型**：分离查询模型

## 总结

这个Go实现展示了如何将C#的DDD和事件驱动架构模式迁移到Go语言。虽然Go缺少一些C#的高级特性（如完整的依赖注入框架），但通过合理的设计，我们仍然可以实现清晰的架构边界和良好的代码组织。 