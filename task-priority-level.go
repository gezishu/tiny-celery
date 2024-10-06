package tinycelery

import "time"

const (
	MinPriorityLevel uint8 = 1
)

func WithPriorityLevel(priorityLevel uint8) taskOption {
	// 优先级作为Unix时间戳秒数
	if priorityLevel < MinPriorityLevel {
		priorityLevel = MinPriorityLevel
	}
	return WithETA(time.Unix(int64(priorityLevel), 0))
}
