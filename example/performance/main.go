package main

import (
	"context"
	"fmt"

	tinycelery "github.com/gezishu/tiny-celery"
	"github.com/gezishu/tiny-celery/example"
)

func main() {
	ctx := context.Background()
	example.RedisClient.FlushAll(ctx)
	client := example.GetTinyCeleryClient()

	tasks := tinycelery.Tasks{}
	for i := 0; i < 1000; i++ {
		tasks = append(tasks, &example.TestBaseTask{
			Desc: fmt.Sprintf("task %d", i),
		})
	}
	example.PanicIfError(client.Delay(
		ctx,
		tasks,
	))
	example.PanicIfError(client.Start(ctx))
}
