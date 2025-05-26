package mediator

import "errors"

var (
	// ErrInvalidRequest 无效的请求错误
	ErrInvalidRequest = errors.New("invalid request type")

	// ErrHandlerNotFound 未找到处理器错误
	ErrHandlerNotFound = errors.New("handler not found")
)
