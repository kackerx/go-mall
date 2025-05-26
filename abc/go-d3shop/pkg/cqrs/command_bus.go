package cqrs

import (
	"context"
	"fmt"
)

// Command 命令接口
type Command interface {
	CommandName() string
}

// CommandHandler 命令处理器接口
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// CommandWithResultHandler 带返回值的命令处理器接口
type CommandWithResultHandler interface {
	Handle(ctx context.Context, cmd Command) (interface{}, error)
}

// CommandMiddleware 命令中间件
type CommandMiddleware func(ctx context.Context, cmd Command, next CommandHandlerFunc) error

// CommandHandlerFunc 命令处理函数
type CommandHandlerFunc func(ctx context.Context, cmd Command) error

// CommandBus 命令总线
type CommandBus struct {
	handlers    map[string]CommandHandler
	middlewares []CommandMiddleware
}

// NewCommandBus 创建命令总线
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers:    make(map[string]CommandHandler),
		middlewares: []CommandMiddleware{},
	}
}

// Register 注册命令处理器
func (b *CommandBus) Register(cmdName string, handler CommandHandler) {
	b.handlers[cmdName] = handler
}

// Use 添加中间件
func (b *CommandBus) Use(middleware CommandMiddleware) {
	b.middlewares = append(b.middlewares, middleware)
}

// Dispatch 分发命令
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	handler, ok := b.handlers[cmd.CommandName()]
	if !ok {
		return fmt.Errorf("no handler registered for command: %s", cmd.CommandName())
	}

	// 构建处理函数
	handlerFunc := func(ctx context.Context, cmd Command) error {
		return handler.Handle(ctx, cmd)
	}

	// 如果没有中间件，直接执行
	if len(b.middlewares) == 0 {
		return handlerFunc(ctx, cmd)
	}

	// 构建中间件链
	var chain CommandHandlerFunc = handlerFunc
	for i := len(b.middlewares) - 1; i >= 0; i-- {
		middleware := b.middlewares[i]
		next := chain
		chain = func(ctx context.Context, cmd Command) error {
			return middleware(ctx, cmd, next)
		}
	}

	return chain(ctx, cmd)
}

// CommandBusWithResult 带返回值的命令总线
type CommandBusWithResult struct {
	handlers    map[string]CommandWithResultHandler
	middlewares []CommandMiddleware
}

// NewCommandBusWithResult 创建带返回值的命令总线
func NewCommandBusWithResult() *CommandBusWithResult {
	return &CommandBusWithResult{
		handlers:    make(map[string]CommandWithResultHandler),
		middlewares: []CommandMiddleware{},
	}
}

// Register 注册命令处理器
func (b *CommandBusWithResult) Register(cmdName string, handler CommandWithResultHandler) {
	b.handlers[cmdName] = handler
}

// Dispatch 分发命令
func (b *CommandBusWithResult) Dispatch(ctx context.Context, cmd Command) (interface{}, error) {
	handler, ok := b.handlers[cmd.CommandName()]
	if !ok {
		return nil, fmt.Errorf("no handler registered for command: %s", cmd.CommandName())
	}

	return handler.Handle(ctx, cmd)
}
