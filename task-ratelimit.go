package tinycelery

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const defaultTaskRateLimit = "1/1s"

// 默认针对单个task限速，仅当多个task需要共享速度限制时，需要指定Key
type taskRateLimit struct {
	Limit    string
	Key      string
	parsed   bool
	max      float64
	duration time.Duration
}

func WithRateLimit(limit string, key string) taskOption {
	return func(m *Meta) error {
		m.RateLimit = taskRateLimit{
			Limit: limit,
			Key:   key,
		}
		return nil
	}
}

func (l *taskRateLimit) limited(ctx context.Context, rc *redis.Client) bool {
	if l.Limit == "" {
		return false
	}
	if !l.parsed {
		if err := l.parse(); err != nil {
			log.Printf("parse rate limit err: %v, switch to default %s\n", err, defaultTaskRateLimit)
			l.Limit = defaultTaskRateLimit
			if err := l.parse(); err != nil {
				panic(fmt.Errorf("unreach err: %v", err))
			}
		}
		l.parsed = true
	}

	now := getNow()
	key := fmt.Sprintf("rate-limit.%s.%d", l.Key, now.UnixMicro()/l.duration.Microseconds())
	defer func() {
		rc.Expire(ctx, key, l.duration*2)
	}()
	res := rc.IncrByFloat(ctx, key, 1.0)
	if err := res.Err(); err != nil {
		log.Printf("check rate limit err: %v\n", err)
		return true
	}
	return res.Val()-l.max > 0.01

}

func (l *taskRateLimit) parse() error {
	ss := strings.Split(l.Limit, "/")
	if len(ss) != 2 {
		return fmt.Errorf("limit %s format invalid", l.Limit)
	}
	max, err := strconv.ParseFloat(ss[0], 64)
	if err != nil {
		return err
	}
	// 目前实现不支持
	if max < 1 {
		return fmt.Errorf("limit %s format unsupport", l.Limit)
	}
	l.max = max
	duration, err := time.ParseDuration(ss[1])
	if err != nil {
		return err
	}
	if duration <= 0 {
		return fmt.Errorf("limit %s format invalid", l.Limit)
	}
	l.duration = duration
	return nil
}
