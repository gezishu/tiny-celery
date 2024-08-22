package tinycelery

import (
	"context"
	"fmt"
	"log"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
)

var (
	marshal   = jsoniter.Marshal
	unmarshal = jsoniter.Unmarshal

	defaultRedisBroker = &redisBroker{
		rc: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
		queue: "tiny-celery-default",
		tasks: map[string]Task{},
	}
)

type redisBroker struct {
	rc    *redis.Client
	queue string
	tasks map[string]Task
}

func (b *redisBroker) send(ctx context.Context, messages ...*Message) error {
	if len(messages) == 0 {
		return nil
	}
	pipline := b.rc.Pipeline()
	for _, message := range messages {
		data, err := marshal(message)
		if err != nil {
			return err
		}
		res := pipline.ZAdd(ctx, b.queue, redis.Z{
			Score:  float64(message.ETA),
			Member: data,
		})
		if err := res.Err(); err != nil {
			return err
		}
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (b *redisBroker) fetch(ctx context.Context) (*Message, error) {
	maxETA := getNow().Unix() + 20
	countRes := b.rc.ZCount(ctx, b.queue, "0", fmt.Sprintf("%d", maxETA))
	if err := countRes.Err(); err != nil {
		return nil, err
	}
	if countRes.Val() == 0 {
		return nil, ErrMessageNotExists
	}
	res := b.rc.ZPopMin(ctx, b.queue, 1)
	if err := res.Err(); err != nil {
		return nil, err
	}
	data := []byte(res.Val()[0].Member.(string))
	meta := &Meta{}
	if err := unmarshal(data, meta); err != nil {
		return nil, err
	}
	task, ok := b.tasks[meta.Name]
	if !ok {
		return nil, ErrTaskNotRegistered
	}
	message := &Message{
		Task: reflect.New(reflect.TypeOf(task).Elem()).Interface().(Task),
	}
	if err := unmarshal(data, message); err != nil {
		return nil, err
	}
	message.setDefault()
	if meta.ETA > maxETA {
		log.Println("task eta faraway from now, restore")
		return nil, b.send(ctx, message)
	}
	return message, nil
}

func (b *redisBroker) getAccumulate(ctx context.Context) (uint32, error) {
	res := b.rc.ZCard(ctx, b.queue)
	if err := res.Err(); err != nil {
		return 0, err
	}
	return uint32(res.Val()), nil
}
