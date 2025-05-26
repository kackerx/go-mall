package behaviors

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/yourusername/go-d3shop/pkg/mediator"
)

// ValidationBehavior 验证管道行为
type ValidationBehavior struct {
	validator *validator.Validate
}

// NewValidationBehavior 创建验证行为
func NewValidationBehavior() *ValidationBehavior {
	return &ValidationBehavior{
		validator: validator.New(),
	}
}

// Handle 处理请求
func (b *ValidationBehavior) Handle(ctx context.Context, request mediator.IRequest, next mediator.RequestHandlerFunc) (interface{}, error) {
	// 验证请求
	err := b.validator.Struct(request)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// 执行下一个处理器
	return next(ctx)
}
