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
)

func IsRedisNilErr(err error) bool {
	return err != nil && errors.Is(err, redis.Nil)
}
