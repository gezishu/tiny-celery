package tinycelery

import (
	"context"
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
		res := pipline.RPush(ctx, b.queue, data)
		if err := res.Err(); err != nil {
			return err
		}
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (b *redisBroker) restore(ctx context.Context, messages ...*Message) error {
	if len(messages) == 0 {
		return nil
	}
	pipline := b.rc.Pipeline()
	for _, message := range messages {
		data, err := marshal(message)
		if err != nil {
			return err
		}
		res := pipline.LPush(ctx, b.queue, data)
		if err := res.Err(); err != nil {
			return err
		}
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (b *redisBroker) fetch(ctx context.Context) (*Message, error) {
	res := b.rc.LPop(ctx, b.queue)
	if err := res.Err(); err != nil {
		if IsRedisNilErr(err) {
			return nil, ErrMessageNotExists
		}
		return nil, err
	}
	data := []byte(res.Val())
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
	return message, nil
}
