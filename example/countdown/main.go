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
	now := time.Now()
	countdown := time.Second * 15
	for i := 0; i < 3; i++ {
		tasks = append(tasks, &example.TestBaseTask{
			Desc: fmt.Sprintf("eta task %d, expect execute at %s", i, now.Add(countdown).Format(time.TimeOnly)),
		})
	}
	example.PanicIfError(client.Delay(
		ctx,
		tasks,
		tinycelery.WithCountdown(countdown),
	))
	example.PanicIfError(client.Start(ctx))
}
