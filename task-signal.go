package tinycelery

import (
	"fmt"
	"time"
)

const (
	PANIC   TaskSignalType = "panic"
	ERROR   TaskSignalType = "error"
	TIMEOUT TaskSignalType = "timeout"
	START   TaskSignalType = "start"
	DONE    TaskSignalType = "done"
)

type TaskSignalType string

type TaskSignal struct {
	Type    TaskSignalType
	Message string
	At      time.Time
}

func (e TaskSignal) String() string {
	if e.Message == "" {
		return string(e.Type)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func newTaskSignal(taskSignalType TaskSignalType, message string) *TaskSignal {
	return &TaskSignal{
		Type:    taskSignalType,
		Message: message,
		At:      getNow(),
	}
}
