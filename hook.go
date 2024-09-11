package tinycelery

import "context"

type TaskHooks struct {
	BeforeCreate  func(context.Context) error
	BeforeExecute func(context.Context) error

	AfterSucceed func(context.Context) error
	AfterFailed  func(context.Context) error
	AfterTimeout func(context.Context) error

	BeforeProcessExit func(context.Context) error
}
