package mediator

import (
	"context"
	"fmt"
)

// HandlerFunc 通用处理函数类型
type HandlerFunc func(ctx context.Context, request interface{}) (interface{}, error)

// NotificationHandlerFunc 通知处理函数类型
type NotificationHandlerFunc func(ctx context.Context, notification interface{}) error

// MediatorV2 更符合Go风格的中介者实现
type MediatorV2 struct {
	handlers             map[string]HandlerFunc
	notificationHandlers map[string][]NotificationHandlerFunc
	pipelineBehaviors    []IPipelineBehavior
}

// NewMediatorV2 创建新的中介者
func NewMediatorV2() *MediatorV2 {
	return &MediatorV2{
		handlers:             make(map[string]HandlerFunc),
		notificationHandlers: make(map[string][]NotificationHandlerFunc),
		pipelineBehaviors:    []IPipelineBehavior{},
	}
}

// RegisterCommandHandler 注册命令处理器
func (m *MediatorV2) RegisterCommandHandler(commandName string, handler HandlerFunc) {
	m.handlers[commandName] = handler
}

// RegisterQueryHandler 注册查询处理器
func (m *MediatorV2) RegisterQueryHandler(queryName string, handler HandlerFunc) {
	m.handlers[queryName] = handler
}

// RegisterNotificationHandler 注册通知处理器
func (m *MediatorV2) RegisterNotificationHandler(eventName string, handler NotificationHandlerFunc) {
	if _, ok := m.notificationHandlers[eventName]; !ok {
		m.notificationHandlers[eventName] = []NotificationHandlerFunc{}
	}
	m.notificationHandlers[eventName] = append(m.notificationHandlers[eventName], handler)
}

// AddPipelineBehavior 添加管道行为
func (m *MediatorV2) AddPipelineBehavior(behavior IPipelineBehavior) {
	m.pipelineBehaviors = append(m.pipelineBehaviors, behavior)
}

// Send 发送请求
func (m *MediatorV2) Send(ctx context.Context, requestName string, request IRequest) (interface{}, error) {
	handler, ok := m.handlers[requestName]
	if !ok {
		return nil, fmt.Errorf("no handler registered for request: %s", requestName)
	}

	// 构建处理函数
	handlerFunc := func(ctx context.Context) (interface{}, error) {
		return handler(ctx, request)
	}

	// 执行管道
	if len(m.pipelineBehaviors) == 0 {
		return handlerFunc(ctx)
	}

	// 构建管道链
	var pipeline RequestHandlerFunc = handlerFunc
	for i := len(m.pipelineBehaviors) - 1; i >= 0; i-- {
		behavior := m.pipelineBehaviors[i]
		next := pipeline
		pipeline = func(ctx context.Context) (interface{}, error) {
			return behavior.Handle(ctx, request, next)
		}
	}

	return pipeline(ctx)
}

// Publish 发布通知
func (m *MediatorV2) Publish(ctx context.Context, eventName string, notification INotification) error {
	handlers, ok := m.notificationHandlers[eventName]
	if !ok {
		return nil // 没有处理器不算错误
	}

	for _, handler := range handlers {
		if err := handler(ctx, notification); err != nil {
			return err
		}
	}

	return nil
}
