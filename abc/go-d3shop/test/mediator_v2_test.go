package test

import (
	"context"
	"testing"

	"github.com/yourusername/go-d3shop/application/commands"
	"github.com/yourusername/go-d3shop/domain/aggregates/order"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// TestMediatorV2CommandHandling 测试MediatorV2的命令处理
func TestMediatorV2CommandHandling(t *testing.T) {
	// 创建MediatorV2
	med := mediator.NewMediatorV2()

	// 注册命令处理器
	med.RegisterCommandHandler("CreateOrder", func(ctx context.Context, request interface{}) (interface{}, error) {
		cmd, ok := request.(*commands.CreateOrderCommandV3)
		if !ok {
			t.Error("类型断言失败")
			return nil, mediator.ErrInvalidRequest
		}

		// 模拟创建订单
		t.Logf("创建订单: Name=%s, Count=%d", cmd.Name, cmd.Count)

		// 返回模拟的订单ID
		return order.NewOrderID(123), nil
	})

	// 创建命令
	cmd := &commands.CreateOrderCommandV3{
		Name:  "测试商品",
		Price: 100,
		Count: 2,
	}

	// 发送命令
	result, err := med.Send(context.Background(), "CreateOrder", cmd)
	if err != nil {
		t.Fatalf("发送命令失败: %v", err)
	}

	// 验证结果
	orderID, ok := result.(order.OrderID)
	if !ok {
		t.Fatal("返回值类型错误")
	}

	if orderID.Value() != 123 {
		t.Errorf("期望订单ID为123，实际为%d", orderID.Value())
	}
}

// TestMediatorV2EventHandling 测试MediatorV2的事件处理
func TestMediatorV2EventHandling(t *testing.T) {
	// 创建MediatorV2
	med := mediator.NewMediatorV2()

	// 计数器，用于验证处理器被调用
	handler1Called := false
	handler2Called := false

	// 注册多个事件处理器
	med.RegisterNotificationHandler("OrderCreated", func(ctx context.Context, notification interface{}) error {
		handler1Called = true
		t.Log("处理器1: 订单创建事件")
		return nil
	})

	med.RegisterNotificationHandler("OrderCreated", func(ctx context.Context, notification interface{}) error {
		handler2Called = true
		t.Log("处理器2: 发送通知邮件")
		return nil
	})

	// 模拟事件
	event := struct {
		OrderID int64
		Name    string
	}{
		OrderID: 123,
		Name:    "测试订单",
	}

	// 发布事件
	err := med.Publish(context.Background(), "OrderCreated", event)
	if err != nil {
		t.Fatalf("发布事件失败: %v", err)
	}

	// 验证两个处理器都被调用
	if !handler1Called {
		t.Error("处理器1未被调用")
	}
	if !handler2Called {
		t.Error("处理器2未被调用")
	}
}

// TestMediatorV2Pipeline 测试MediatorV2的管道行为
func TestMediatorV2Pipeline(t *testing.T) {
	// 创建MediatorV2
	med := mediator.NewMediatorV2()

	// 执行顺序记录
	var executionOrder []string

	// 添加管道行为
	med.AddPipelineBehavior(&TestPipelineBehavior{
		name:  "Behavior1",
		order: &executionOrder,
	})

	med.AddPipelineBehavior(&TestPipelineBehavior{
		name:  "Behavior2",
		order: &executionOrder,
	})

	// 注册命令处理器
	med.RegisterCommandHandler("TestCommand", func(ctx context.Context, request interface{}) (interface{}, error) {
		executionOrder = append(executionOrder, "Handler")
		return "OK", nil
	})

	// 发送命令
	_, err := med.Send(context.Background(), "TestCommand", struct{}{})
	if err != nil {
		t.Fatalf("发送命令失败: %v", err)
	}

	// 验证执行顺序
	expectedOrder := []string{
		"Behavior1-Before",
		"Behavior2-Before",
		"Handler",
		"Behavior2-After",
		"Behavior1-After",
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("执行顺序长度不匹配: 期望%d，实际%d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("执行顺序[%d]不匹配: 期望%s，实际%s", i, expected, executionOrder[i])
		}
	}
}

// TestPipelineBehavior 测试用的管道行为
type TestPipelineBehavior struct {
	name  string
	order *[]string
}

func (b *TestPipelineBehavior) Handle(ctx context.Context, request mediator.IRequest, next mediator.RequestHandlerFunc) (interface{}, error) {
	*b.order = append(*b.order, b.name+"-Before")
	result, err := next(ctx)
	*b.order = append(*b.order, b.name+"-After")
	return result, err
}
