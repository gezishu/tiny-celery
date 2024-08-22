package example

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	tinycelery "github.com/gezishu/tiny-celery"
)

type TestHookTask struct {
	Desc  string
	Sleep time.Duration
	Error string
	Panic string
}

func (t *TestHookTask) Hooks(ctx context.Context) tinycelery.TaskHooks {
	return tinycelery.TaskHooks{
		tinycelery.BeforeCreate: func(ctx context.Context) error {
			log.Println(t, tinycelery.BeforeCreate)
			return nil
		},
		tinycelery.BeforeExecute: func(ctx context.Context) error {
			log.Println(t, tinycelery.BeforeExecute)
			return nil
		},
		tinycelery.AfterSucceed: func(ctx context.Context) error {
			log.Println(t, tinycelery.AfterSucceed)
			return nil
		},
		tinycelery.AfterFailed: func(ctx context.Context) error {
			log.Println(t, tinycelery.AfterFailed)
			return nil
		},
		tinycelery.AfterTimeout: func(ctx context.Context) error {
			fmt.Println(t, tinycelery.AfterTimeout)
			return nil
		},
		tinycelery.BeforeProcessExit: func(ctx context.Context) error {
			fmt.Println(t, tinycelery.BeforeProcessExit)
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
