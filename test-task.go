package tinycelery

import (
	"context"
	"fmt"
	"time"
)

type TestTask struct {
	N uint64
}

func (t *TestTask) Hooks(ctx context.Context) TaskHooks {
	return TaskHooks{
		// BeforeCreate: func(ctx context.Context) error {
		// 	fmt.Println(t.N, BeforeCreate)
		// 	return nil
		// },
		// BeforeExecute: func(ctx context.Context) error {
		// 	fmt.Println(t.N, BeforeExecute)
		// 	return nil
		// },
		// AfterSucceed: func(ctx context.Context) error {
		// 	fmt.Println(t.N, AfterSucceed)
		// 	return nil
		// },
		// AfterFailed: func(ctx context.Context) error {
		// 	fmt.Println(t.N, AfterFailed)
		// 	return nil
		// },
		// AfterTimeout: func(ctx context.Context) error {
		// 	fmt.Println(t.N, AfterTimeout)
		// 	return nil
		// },
		// BeforeProcessExit: func(ctx context.Context) error {
		// 	fmt.Println(t.N, BeforeProcessExit)
		// 	return nil
		// },
	}
}
func (t *TestTask) Execute(ctx context.Context) error {
	fmt.Println(getNow(), "test task start", t.N)
	time.Sleep(time.Second * 0)
	// fmt.Println(GetNow(), "test task done", t.N)
	return nil
}
