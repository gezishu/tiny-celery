package example

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	tinycelery "github.com/gezishu/tiny-celery"
)

func init() {
	tasks = append(tasks, &TestHookTask{})
}

type TestHookTask struct {
	Desc  string
	Sleep time.Duration
	Error string
	Panic string
}

func (t *TestHookTask) Hooks(ctx context.Context) *tinycelery.TaskHooks {
	return &tinycelery.TaskHooks{
		BeforeCreate: func(ctx context.Context) error {
			log.Println(t, "before-create")
			return nil
		},
		BeforeExecute: func(ctx context.Context) error {
			log.Println(t, "before-execute")
			return nil
		},
		AfterSucceed: func(ctx context.Context) error {
			log.Println(t, "after-succeed")
			return nil
		},
		AfterFailed: func(ctx context.Context) error {
			log.Println(t, "after-failed")
			return nil
		},
		AfterTimeout: func(ctx context.Context) error {
			fmt.Println(t, "after-timeout")
			return nil
		},
		BeforeProcessExit: func(ctx context.Context) error {
			fmt.Println(t, "before-process-exit")
			return nil
		},
	}
}
func (t *TestHookTask) Execute(ctx context.Context) error {
	log.Println(t, "start")
	if t.Panic != "" {
		panic(t.Panic)
	} else if t.Error != "" {
		return errors.New(t.Error)
	}
	time.Sleep(t.Sleep)
	log.Println(t, "succeed")
	return nil
}
