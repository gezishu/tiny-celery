package tinycelery

import "time"

const (
	minTaskTimeLimit     = time.Second * 1
	defaultTaskTimeLimit = time.Minute * 2
)

func WithTimeLimit(timeLimit time.Duration) taskOption {
	return func(m *Meta) error {
		m.TimeLimit = timeLimit
		return nil
	}
}
