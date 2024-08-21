package tinycelery

import "time"

type TaskOption func(m *Meta) error

func WithTimeLimit(timeLimit time.Duration) TaskOption {
	return func(m *Meta) error {
		m.TaskTimeLimit = timeLimit
		return nil
	}
}

func WithRateLimit(rateLimit TaskRateLimit) TaskOption {
	return func(m *Meta) error {
		m.TaskRateLimit = rateLimit
		return nil
	}
}
