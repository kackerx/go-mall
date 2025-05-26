package command

import "context"

type Command interface {
	CommandName() string
}

type CommandHandler interface {
	Handle(ctx context.Context, command Command) error
}
