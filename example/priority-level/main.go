package main

import (
	"context"
	"time"

	tinycelery "github.com/gezishu/tiny-celery"
	"github.com/gezishu/tiny-celery/example"
)

func main() {
	ctx := context.Background()
	example.RedisClient.FlushAll(ctx)
	client := example.GetTinyCeleryClient()
	tasks := tinycelery.Tasks{}
	for i := 0; i < 20; i++ {
		tasks = append(tasks, &example.TestBaseTask{
			Desc:  "without priority level",
			Sleep: time.Second * 5,
		})
	}
	example.PanicIfError(client.Delay(
		ctx,
		tasks,
	))
	tasksWithPriority := tinycelery.Tasks{}
	for i := 0; i < 20; i++ {
		tasksWithPriority = append(tasksWithPriority, &example.TestBaseTask{
			Desc:  "with priority level",
			Sleep: time.Second * 5,
		})
	}
	example.PanicIfError(client.Delay(
		ctx,
		tasksWithPriority,
		tinycelery.WithPriorityLevel(5),
	))
	example.PanicIfError(client.Start(ctx))
}
