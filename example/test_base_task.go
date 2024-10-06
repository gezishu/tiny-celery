package example

import (
	"context"
	"log"
	"time"

	tinycelery "github.com/gezishu/tiny-celery"
)

func init() {
	tasks = append(tasks, &TestBaseTask{})
}

type TestBaseTask struct {
	Desc  string
	Sleep time.Duration
}

func (t *TestBaseTask) Hooks(ctx context.Context) *tinycelery.TaskHooks {
	return nil
}
func (t *TestBaseTask) Execute(ctx context.Context) error {
	log.Println(t, "start")
	time.Sleep(t.Sleep)
	return nil
}
