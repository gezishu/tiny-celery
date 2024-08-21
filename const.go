package tinycelery

import "time"

const (
	defaultTaskTimeLimit = time.Minute * 2
	defaultTaskRateLimit = "1/1s"
	defaultCurrency      = 1
	defaultPrefetch      = 0
)
