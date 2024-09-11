package tinycelery

import "context"

type Task interface {
	Hooks(context.Context) *TaskHooks
	Execute(context.Context) error
}

type Tasks []Task
