package mediator

import (
	"context"
	"fmt"
	"reflect"
)

// IRequest 请求接口
type IRequest interface{}

// ICommand 命令接口
type ICommand interface {
	IRequest
}

// IQuery 查询接口
type IQuery interface {
	IRequest
}

// INotification 通知接口（用于事件）
type INotification interface{}

// IRequestHandler 请求处理器接口
type IRequestHandler[TRequest IRequest, TResponse any] interface {
	Handle(ctx context.Context, request TRequest) (TResponse, error)
}

// ICommandHandler 命令处理器接口
type ICommandHandler[TCommand ICommand] interface {
	Handle(ctx context.Context, command TCommand) error
}

// ICommandWithResultHandler 带返回值的命令处理器接口
type ICommandWithResultHandler[TCommand ICommand, TResult any] interface {
	Handle(ctx context.Context, command TCommand) (TResult, error)
}

// IQueryHandler 查询处理器接口
type IQueryHandler[TQuery IQuery, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

// INotificationHandler 通知处理器接口
type INotificationHandler[TNotification INotification] interface {
	Handle(ctx context.Context, notification TNotification) error
}

// IPipelineBehavior 管道行为接口
type IPipelineBehavior interface {
	Handle(ctx context.Context, request IRequest, next RequestHandlerFunc) (interface{}, error)
}

// RequestHandlerFunc 请求处理函数
type RequestHandlerFunc func(ctx context.Context) (interface{}, error)

// IMediator 中介者接口
type IMediator interface {
	Send(ctx context.Context, request IRequest) (interface{}, error)
	Publish(ctx context.Context, notification INotification) error
}

// Mediator 中介者实现
type Mediator struct {
	handlers             map[reflect.Type]interface{}
	notificationHandlers map[reflect.Type][]interface{}
	pipelineBehaviors    []IPipelineBehavior
}

// NewMediator 创建新的中介者
func NewMediator() *Mediator {
	return &Mediator{
		handlers:             make(map[reflect.Type]interface{}),
		notificationHandlers: make(map[reflect.Type][]interface{}),
		pipelineBehaviors:    []IPipelineBehavior{},
	}
}

// RegisterHandler 注册处理器
func (m *Mediator) RegisterHandler(requestType reflect.Type, handler interface{}) {
	m.handlers[requestType] = handler
}

// RegisterNotificationHandler 注册通知处理器
func (m *Mediator) RegisterNotificationHandler(notificationType reflect.Type, handler interface{}) {
	if _, ok := m.notificationHandlers[notificationType]; !ok {
		m.notificationHandlers[notificationType] = []interface{}{}
	}
	m.notificationHandlers[notificationType] = append(m.notificationHandlers[notificationType], handler)
}

// AddPipelineBehavior 添加管道行为
func (m *Mediator) AddPipelineBehavior(behavior IPipelineBehavior) {
	m.pipelineBehaviors = append(m.pipelineBehaviors, behavior)
}

// Send 发送请求
func (m *Mediator) Send(ctx context.Context, request IRequest) (interface{}, error) {
	requestType := reflect.TypeOf(request)
	handler, ok := m.handlers[requestType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for request type %v", requestType)
	}

	// 构建处理函数
	handlerFunc := func(ctx context.Context) (interface{}, error) {
		handlerValue := reflect.ValueOf(handler)
		method := handlerValue.MethodByName("Handle")
		if !method.IsValid() {
			return nil, fmt.Errorf("handler does not have Handle method")
		}

		args := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(request),
		}

		results := method.Call(args)
		if len(results) == 2 {
			// 有返回值的情况
			if !results[1].IsNil() {
				return results[0].Interface(), results[1].Interface().(error)
			}
			return results[0].Interface(), nil
		} else if len(results) == 1 {
			// 只有error返回值的情况
			if !results[0].IsNil() {
				return nil, results[0].Interface().(error)
			}
			return nil, nil
		}

		return nil, fmt.Errorf("unexpected number of return values")
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
func (m *Mediator) Publish(ctx context.Context, notification INotification) error {
	notificationType := reflect.TypeOf(notification)
	handlers, ok := m.notificationHandlers[notificationType]
	if !ok {
		return nil // 没有处理器不算错误
	}

	for _, handler := range handlers {
		handlerValue := reflect.ValueOf(handler)
		method := handlerValue.MethodByName("Handle")
		if !method.IsValid() {
			continue
		}

		args := []reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(notification),
		}

		results := method.Call(args)
		if len(results) == 1 && !results[0].IsNil() {
			if err, ok := results[0].Interface().(error); ok && err != nil {
				return err
			}
		}
	}

	return nil
}
