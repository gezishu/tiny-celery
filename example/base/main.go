package main

import (
	"context"

	tinycelery "github.com/gezishu/tiny-celery"
	"github.com/gezishu/tiny-celery/example"
)

func main() {
	ctx := context.Background()
	example.RedisClient.FlushAll(ctx)
	client := example.GetTinyCeleryClient()

	example.PanicIfError(client.Delay(
		ctx,
		tinycelery.Tasks{
			&example.TestBaseTask{
				Desc: "test base",
			},
		},
	))

	example.PanicIfError(client.Start(ctx))
}
