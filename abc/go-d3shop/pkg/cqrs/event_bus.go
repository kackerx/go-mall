package cqrs

import (
	"context"
	"sync"
)

// Event 事件接口
type Event interface {
	EventName() string
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

// EventHandlerFunc 事件处理函数
type EventHandlerFunc func(ctx context.Context, event Event) error

// Handle 实现EventHandler接口
func (f EventHandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// EventBus 事件总线
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe 订阅事件
func (b *EventBus) Subscribe(eventName string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.handlers[eventName]; !ok {
		b.handlers[eventName] = []EventHandler{}
	}
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

// SubscribeFunc 使用函数订阅事件
func (b *EventBus) SubscribeFunc(eventName string, handler EventHandlerFunc) {
	b.Subscribe(eventName, handler)
}

// Publish 发布事件
func (b *EventBus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	handlers, ok := b.handlers[event.EventName()]
	b.mu.RUnlock()

	if !ok {
		return nil // 没有处理器不是错误
	}

	// 复制处理器列表，避免在处理过程中被修改
	handlersCopy := make([]EventHandler, len(handlers))
	copy(handlersCopy, handlers)

	// 同步执行所有处理器
	for _, handler := range handlersCopy {
		if err := handler.Handle(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

// PublishAsync 异步发布事件
func (b *EventBus) PublishAsync(ctx context.Context, event Event) {
	b.mu.RLock()
	handlers, ok := b.handlers[event.EventName()]
	b.mu.RUnlock()

	if !ok {
		return
	}

	// 复制处理器列表
	handlersCopy := make([]EventHandler, len(handlers))
	copy(handlersCopy, handlers)

	// 异步执行每个处理器
	for _, handler := range handlersCopy {
		go func(h EventHandler) {
			_ = h.Handle(ctx, event)
		}(handler)
	}
}
