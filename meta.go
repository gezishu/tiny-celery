package tinycelery

import (
	"fmt"
	"sync"
	"time"
)

type Meta struct {
	mu        sync.Mutex
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	TimeLimit time.Duration `json:"timelimit"`
	RateLimit taskRateLimit `json:"ratelimit"`
	state     taskState
	rtName    string
}

type taskOption func(m *Meta) error

func (m *Meta) setState(targetState taskState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	currState := m.state
	for _, state := range taskStateTransformMap[currState] {
		if state == targetState {
			m.state = targetState
			return nil
		}
	}
	return fmt.Errorf("can't set task state from %v to %v", currState, targetState)
}

func (m *Meta) setDefault() {
	if m.ID == "" {
		m.ID = genRandString(8)
	}
	if m.TimeLimit < minTaskTimeLimit {
		m.TimeLimit = defaultTaskTimeLimit
	}
	m.state = taskINIT
}
