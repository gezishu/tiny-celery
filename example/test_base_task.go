package example

import (
	"context"
	"log"

	tinycelery "github.com/gezishu/tiny-celery"
)

type TestBaseTask struct {
	Desc string
}

func (t *TestBaseTask) Hooks(ctx context.Context) *tinycelery.TaskHooks {
	return nil
}
func (t *TestBaseTask) Execute(ctx context.Context) error {
	log.Println(t, "start")
	return nil
}
