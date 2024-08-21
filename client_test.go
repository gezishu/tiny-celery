package tinycelery

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestExecuteTask(t *testing.T) {
	ctx := context.Background()
	rc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	rc.FlushDB(ctx)
	client, err := NewClient().Init(
		rc,
		WithCurrencyMax(4),
		WithPrefetch(8),
		WithQueue("test"),
	)
	if err != nil {
		panic(err)
	}
	err = client.RegisterTask(&TestTask{})
	if err != nil {
		panic(err)
	}
	tasks := []TaskInterface{}
	for i := 0; i < 40; i++ {
		tasks = append(tasks, &TestTask{
			N: uint64(i),
		})
	}
	err = client.Delay(
		ctx,
		tasks,
		// WithTimeout(time.Second*2),
		WithRateLimit(TaskRateLimit{
			Limit: "1/2s",
		}),
	)
	if err != nil {
		panic(err)
	}
	if err := client.Start(ctx); err != nil {
		panic(err)
	}

}
