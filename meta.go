package tinycelery

import (
	"fmt"
	"sync"
	"time"
)

const (
	StateINIT TaskState = iota
	StateRUNNING
	StateSUCCEED
	StateFAILED
	StateTIMEOUT
)

var stateTransformMap = map[TaskState][]TaskState{
	StateINIT:    {StateRUNNING},
	StateRUNNING: {StateSUCCEED, StateFAILED, StateTIMEOUT},
	StateFAILED:  {StateRUNNING},
	StateTIMEOUT: {StateRUNNING},
}

type TaskState uint8

type Meta struct {
	TaskID        string
	TaskName      string
	TaskTimeLimit time.Duration
	TaskRateLimit TaskRateLimit
	mu            sync.Mutex `json:"-"`
	taskState     TaskState  `json:"-"`
	taskRTName    string     `json:"-"`
}

func (m *Meta) SetTaskState(taskState TaskState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	currState := m.taskState
	for _, state := range stateTransformMap[currState] {
		if state == taskState {
			m.taskState = taskState
			return nil
		}
	}
	return fmt.Errorf("can't set task state from %v to %v", currState, taskState)
}

func (m *Meta) SetDefault() {
	if m.TaskTimeLimit < time.Second {
		m.TaskTimeLimit = defaultTaskTimeLimit
	}
	if m.TaskRateLimit.Limit != "" {
		if m.TaskRateLimit.Key == "" {
			m.TaskRateLimit.Key = m.TaskName
		}
	}
	m.taskState = StateINIT
	m.taskRTName = fmt.Sprintf("%s[%s]", m.TaskName, genRandString(4))
}
