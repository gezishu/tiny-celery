package main

import (
	"context"
	"fmt"
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
		tasks = append(tasks, &example.TestHookTask{
			Desc:  fmt.Sprintf("task %d", i),
			Sleep: time.Second * 30,
		})
	}
	example.PanicIfError(client.Delay(
		ctx,
		tasks,
	))
	example.PanicIfError(client.Start(ctx))
}
