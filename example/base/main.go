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

	example.PanicIfError(client.Delay(
		ctx,
		tinycelery.Tasks{
			&example.TestBaseTask{
				Desc: "test succeed",
			},
		},
	))

	example.PanicIfError(client.Delay(
		ctx,
		tinycelery.Tasks{
			&example.TestBaseTask{
				Desc:  "test timeout",
				Sleep: time.Second * 5,
			},
		},
		tinycelery.WithTimeLimit(time.Second*2),
	))

	example.PanicIfError(client.Delay(
		ctx,
		tinycelery.Tasks{
			&example.TestBaseTask{
				Desc:  "test failed",
				Error: "test error",
			},
		},
	))

	example.PanicIfError(client.Delay(
		ctx,
		tinycelery.Tasks{
			&example.TestBaseTask{
				Desc:  "test panic",
				Panic: "test panic",
			},
		},
	))

	example.PanicIfError(client.Start(ctx))
}
