package tinycelery

import "time"

func WithETA(eta time.Time) taskOption {
	return func(m *Meta) error {
		m.ETA = eta.Unix()
		return nil
	}
}

func WithCountdown(countdown time.Duration) taskOption {
	return WithETA(getNow().Add(countdown))
}
