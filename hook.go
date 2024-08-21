package tinycelery

import "context"

var (
	BeforeCreate  TaskHook = "before-create"
	BeforeExecute TaskHook = "before-execute"

	AfterSucceed TaskHook = "after-succeed"
	AfterFailed  TaskHook = "after-failed"
	AfterTimeout TaskHook = "after-timeout"

	BeforeProcessExit TaskHook = "before-process-exit"
)

type TaskHook string

type TaskHookFunc func(context.Context) error

type TaskHooks map[TaskHook]TaskHookFunc
