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
)

type RedisBroker struct {
	rc    *redis.Client
	queue string
	tasks map[string]TaskInterface
}

func (b *RedisBroker) send(ctx context.Context, messages ...*Message) error {
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

func (b *RedisBroker) restore(ctx context.Context, messages ...*Message) error {
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

func (b *RedisBroker) fetch(ctx context.Context) (*Message, error) {
	res := b.rc.LPop(ctx, b.queue)
	if err := res.Err(); err != nil {
		if IsRedisNilErr(err) {
			return nil, ErrMessageNotExists
		}
		return nil, err
	}
	data := []byte(res.Val())
	meta := Meta{}
	if err := unmarshal(data, &meta); err != nil {
		return nil, err
	}
	taskName := meta.TaskName
	task, ok := b.tasks[taskName]
	if !ok {
		return nil, ErrTaskNotRegistered
	}
	message := &Message{
		Meta: meta,
		Task: reflect.New(reflect.TypeOf(task).Elem()).Interface().(TaskInterface),
	}
	if err := unmarshal(data, message); err != nil {
		return nil, err
	}
	message.SetDefault()
	return message, nil
}
