package tinycelery

import (
	"context"
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

var (
	onTaskStart  = func(ctx context.Context, message *Message) {}
	onTaskDone   = func(ctx context.Context, message *Message) {}
	onTaskFailed = func(_ context.Context, message *Message) {
		logger.Printf("task %s%v failed, err: %v\n", message.rtName, message.Task, message.err)
	}
	onTaskTimeout = func(_ context.Context, message *Message) {
		logger.Printf("task %s%v timeout\n", message.rtName, message.Task)
	}
	onProcessExit = func(_ context.Context, message *Message) {
		logger.Printf("task %s%v on process exit\n", message.rtName, message.Task)
	}
)

func SetOnTaskStart(f func(ctx context.Context, message *Message)) {
	onTaskStart = f
}

func SetOnTaskDone(f func(ctx context.Context, message *Message)) {
	onTaskDone = f
}

func SetOnTaskFailed(f func(ctx context.Context, message *Message)) {
	onTaskFailed = f
}

func SetOnTaskTimeout(f func(ctx context.Context, message *Message)) {
	onTaskTimeout = f
}

func SetOnProcessExit(f func(ctx context.Context, message *Message)) {
	onProcessExit = f
}
