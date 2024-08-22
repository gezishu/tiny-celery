package tinycelery

const (
	taskINIT taskState = iota + 1
	taskRUNNING
	taskSUCCEED
	taskFAILED
	taskTIMEOUT
)

var taskStateTransformMap = map[taskState][]taskState{
	taskINIT:    {taskRUNNING},
	taskRUNNING: {taskSUCCEED, taskFAILED, taskTIMEOUT},
	taskFAILED:  {taskRUNNING},
	taskTIMEOUT: {taskRUNNING},
}

type taskState uint8
