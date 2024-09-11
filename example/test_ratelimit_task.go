package example

import (
	"context"
	"log"

	tinycelery "github.com/gezishu/tiny-celery"
)

type TestRateLimitTask struct {
	Desc string
}

func (t *TestRateLimitTask) Hooks(ctx context.Context) *tinycelery.TaskHooks {
	return nil
}

func (t *TestRateLimitTask) Execute(ctx context.Context) error {
	log.Println(t, "start")
	return nil
}
