package cqrs

import "errors"

var (
	// ErrInvalidCommand 无效的命令错误
	ErrInvalidCommand = errors.New("invalid command type")

	// ErrHandlerNotFound 未找到处理器错误
	ErrHandlerNotFound = errors.New("handler not found")

	// ErrInvalidEvent 无效的事件错误
	ErrInvalidEvent = errors.New("invalid event type")
)
