package tinycelery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// 默认针对单个task限速，仅当多个task需要共享速度限制时，需要指定Key
type TaskRateLimit struct {
	Limit    string
	Key      string
	parsed   bool
	max      float64
	duration time.Duration
}

func (l *TaskRateLimit) enable() bool {
	return l.Limit != ""
}

func (l *TaskRateLimit) parse() {
	if l.parsed {
		return
	}
	if err := l._parse(); err != nil {
		log.Printf("parse rate limit failed: %v, switch to default %s\n", err, defaultTaskRateLimit)
		l.Limit = defaultTaskRateLimit
		if err := l._parse(); err != nil {
			panic(err)
		}
	}
	l.parsed = true
}

func (l *TaskRateLimit) _parse() error {
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

func (l *TaskRateLimit) check(ctx context.Context, rc *redis.Client) (bool, error) {
	if !l.parsed {
		return true, errors.New("parse required")
	}
	now := getNow()
	key := fmt.Sprintf("rate-limit.%s.%d", l.Key, now.UnixMicro()/l.duration.Microseconds())
	defer func() {
		rc.Expire(ctx, key, l.duration*2)
	}()
	res := rc.IncrByFloat(ctx, key, 1.0)
	if err := res.Err(); err != nil {
		return true, err
	}
	return res.Val()-l.max > 0.01, nil
}
