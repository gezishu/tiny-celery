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
	for i := 0; i < 50; i++ {
		tasks = append(tasks, &example.TestRateLimitTask{
			Desc: fmt.Sprintf("task %d", i),
		})
	}
	example.PanicIfError(client.Delay(
		ctx,
		tasks,
		tinycelery.WithRateLimit("1/5s", ""),
	))
	example.PanicIfError(client.Start(ctx))
}
