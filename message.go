package tinycelery

import (
	"context"
)

type Message struct {
	Meta
	Task TaskInterface
}

type TaskInterface interface {
	Hooks(context.Context) TaskHooks
	Execute(context.Context) error
}
