# 领域事件 vs 集成事件对比

## 领域事件（Domain Events）

### 特征
- **范围**：限于单个限界上下文内
- **处理**：同步处理，在同一事务内
- **一致性**：强一致性
- **耦合**：进程内耦合
- **失败处理**：失败时整个事务回滚

### 使用场景
1. **订单创建后自动创建发货记录**
   ```go
   OrderCreatedDomainEvent → 创建DeliverRecord（同一事务）
   ```

2. **用户注册后初始化用户配置**
   ```go
   UserRegisteredDomainEvent → 创建UserProfile（必须成功）
   ```

3. **库存扣减后更新统计信息**
   ```go
   InventoryDeductedDomainEvent → 更新统计表（强一致）
   ```

### 实现示例
```go
// 领域事件在聚合根中触发
func (order *Order) Complete() {
    order.Status = "completed"
    order.AddDomainEvent(NewOrderCompletedDomainEvent(order))
}

// 在同一事务中处理
tx.Begin()
orderRepo.Save(order)           // 保存订单
statsRepo.UpdateStats(order)    // 更新统计（领域事件处理）
tx.Commit()                     // 全部成功或全部失败
```

## 集成事件（Integration Events）

### 特征
- **范围**：跨限界上下文/微服务
- **处理**：异步处理，通过消息队列
- **一致性**：最终一致性
- **耦合**：松耦合
- **失败处理**：重试机制，补偿事务

### 使用场景
1. **订单支付后通知库存服务**
   ```go
   OrderPaidIntegrationEvent → 库存服务扣减库存（可能延迟）
   ```

2. **用户注册后发送欢迎邮件**
   ```go
   UserRegisteredIntegrationEvent → 邮件服务发送邮件（允许失败）
   ```

3. **订单完成后更新用户积分**
   ```go
   OrderCompletedIntegrationEvent → 积分服务增加积分（最终一致）
   ```

### 实现示例
```go
// 集成事件显式发布
orderService.CompleteOrder(orderId)                    // 完成订单
eventPublisher.Publish(OrderCompletedIntegrationEvent) // 发布事件

// 其他服务异步处理
// 积分服务
func HandleOrderCompleted(event) {
    retry(3) {
        userService.AddPoints(event.UserId, event.Points)
    }
}
```

## 选择指南

### 使用领域事件当：
- ✅ 需要强一致性
- ✅ 在同一个限界上下文内
- ✅ 操作必须全部成功或全部失败
- ✅ 处理逻辑简单且快速

### 使用集成事件当：
- ✅ 可以接受最终一致性
- ✅ 跨服务/限界上下文通信
- ✅ 处理可能耗时较长
- ✅ 需要解耦服务依赖
- ✅ 允许部分失败和重试

## 混合使用示例

```go
// 1. 用户下单
func CreateOrder(order Order) {
    // 保存订单
    orderRepo.Save(order)
    
    // 领域事件：同步创建发货记录（必须成功）
    deliveryRepo.CreateDelivery(order.Id)
    
    // 集成事件：异步通知其他服务
    eventBus.Publish(OrderCreatedIntegrationEvent{
        OrderId: order.Id,
        UserId: order.UserId,
    })
}

// 2. 库存服务收到事件（可能延迟）
func HandleOrderCreated(event OrderCreatedIntegrationEvent) {
    // 扣减库存（允许失败重试）
    inventory.Reserve(event.OrderId)
}
``` 