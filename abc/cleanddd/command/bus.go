package command

import "context"

type CommandBus interface {
	Register(cmd Command, handler CommandHandler)
	Send(ctx context.Context, cmd Command) error
}

type commandBus struct {
	handlers map[string]CommandHandler
}

func NewCommandBus() CommandBus {
	return &commandBus{
		handlers: make(map[string]CommandHandler),
	}
}

func (b *commandBus) Register(cmd Command, handler CommandHandler) {
	b.handlers[cmd.CommandName()] = handler
}

func (b *commandBus) Send(ctx context.Context, cmd Command) error {
	return b.handlers[cmd.CommandName()].Handle(ctx, cmd)
}
