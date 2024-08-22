package example

import (
	tinycelery "github.com/gezishu/tiny-celery"
	"github.com/redis/go-redis/v9"
)

var tasks = tinycelery.Tasks{
	&TestBaseTask{},
	&TestRateLimitTask{},
}

var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0,
})

func GetTinyCeleryClient() *tinycelery.Client {
	client, err := tinycelery.NewClient().Init(
		RedisClient,
		tinycelery.WithConcurrency(4),
		tinycelery.WithPrefetch(4),
		tinycelery.WithQueue("test"),
	)
	PanicIfError(err)
	for _, task := range tasks {
		PanicIfError(client.RegisterTask(task))
	}
	return client
}

func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
