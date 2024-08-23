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
	example.PanicIfError(client.Delay(ctx, tinycelery.Tasks{
		&example.TestBaseTask{},
		&example.TestHookTask{},
	}))
	accumulate, err := client.GetAccumulate(ctx)
	example.PanicIfError(err)
	fmt.Println("client.GetAccumulate", accumulate)
}
