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
	eta := time.Now().Add(time.Second * 10)
	for i := 0; i < 20; i++ {
		tasks = append(tasks, &example.TestBaseTask{
			Desc: fmt.Sprintf("eta task %d, expect execute at %s", i, eta.Format(time.TimeOnly)),
		})
	}
	example.PanicIfError(client.Delay(
		ctx,
		tasks,
		tinycelery.WithETA(eta),
	))
	example.PanicIfError(client.Start(ctx))
}
