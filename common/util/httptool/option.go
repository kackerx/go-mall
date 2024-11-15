package httptool

import (
	"context"
	"time"
)

type requestOption struct {
	ctx     context.Context
	timeout time.Duration
	data    []byte
	headers map[string]string
}

func defaultRequestOption() *requestOption {
	return &requestOption{
		ctx:     context.Background(),
		timeout: 5 * time.Second,
		data:    nil,
		headers: make(map[string]string),
	}
}

type Option func(opt *requestOption) error

func WithContext(ctx context.Context) Option {
	return func(opt *requestOption) error {
		opt.ctx = ctx
		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opt *requestOption) error {
		opt.timeout = timeout
		return nil
	}
}

func WithData(data []byte) Option {
	return func(opt *requestOption) error {
		opt.data = data
		return nil
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(opt *requestOption) error {
		for k, v := range headers {
			opt.headers[k] = v
		}
		return nil
	}
}
