package tinycelery

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrTaskHasRegistered = errors.New("task has registered")
	ErrTaskNotRegistered = errors.New("task not registered")
	ErrMessageIsNil      = errors.New("message is nil")
	ErrMessageNotExists  = errors.New("message not exists")
	ErrNoMatchedTask     = errors.New("no matched task")
)

func IsRedisNilErr(err error) bool {
	return err != nil && errors.Is(err, redis.Nil)
}
