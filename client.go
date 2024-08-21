package tinycelery

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	sync.Mutex
	currency Currency
	prefetch uint32
	messages []*Message
	broker   *RedisBroker
	logger   *log.Logger
	state    ClientState
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) setDefault() {
	c.currency = Currency{
		max:  defaultCurrency,
		curr: 0,
		C:    make(chan uint32, 1),
	}
	c.prefetch = defaultPrefetch
	c.logger = log.New(os.Stdout, "[tiny-celery] ", log.Ltime)
}

func (c *Client) setState(state ClientState) {
	c.Lock()
	defer c.Unlock()
	c.state = state
}

func (c *Client) Init(rc *redis.Client, initOptions ...initOption) (*Client, error) {
	c.setDefault()
	c.broker = &RedisBroker{
		rc:    rc,
		queue: "tiny-celery-default",
		tasks: map[string]TaskInterface{},
	}
	for _, initOption := range initOptions {
		err := initOption(c)
		if err != nil {
			return nil, err
		}
	}
	c.messages = make([]*Message, c.currency.max+c.prefetch)
	return c, nil
}

func (c *Client) RegisterTask(task TaskInterface) error {
	taskName := getType(task)
	if _, registered := c.broker.tasks[taskName]; registered {
		return ErrTaskHasRegistered
	}
	c.broker.tasks[taskName] = task
	return nil
}

func (c *Client) updateMessageByIndex(ctx context.Context, index int) error {
	c.messages[index] = nil
	message, err := c.broker.fetch(ctx)
	if err != nil {
		return err
	}
	c.messages[index] = message
	return nil
}

func (c *Client) scanMessages(ctx context.Context) {
	if c.state != ClientRUNNING {
		return
	}
	for i := 0; i < len(c.messages); i++ {
		if c.messages[i] == nil {
			if err := c.updateMessageByIndex(ctx, i); err != nil {
				if errors.Is(err, ErrMessageNotExists) {
					break
				} else {
					c.logger.Printf("update message by index err: %v\n", err)
				}
				continue
			}
		}
	}
	for i := 0; i < len(c.messages); i++ {
		message := c.messages[i]
		if message == nil {
			continue
		}
		switch message.taskState {
		case StateINIT:
			var runable = true
			if !c.currency.prerequire() {
				runable = false
			}
			if !runable {
				continue
			}
			if message.TaskRateLimit.enable() {
				message.TaskRateLimit.parse()
				limited, err := message.TaskRateLimit.check(ctx, c.broker.rc)
				if err != nil {
					runable = false
				} else if limited {
					runable = false
				}
			}
			if !runable {
				continue
			}
			if c.currency.require() {
				if err := message.SetTaskState(StateRUNNING); err != nil {
					// 理论上不可达
					panic(err)
				}
				go c.runTask(ctx, message)
			}
		case StateSUCCEED:
			c.updateMessageByIndex(ctx, i)
		case StateFAILED:
			// todo 增加重试逻辑
			c.logger.Printf("task %s failed\n", message.taskRTName)
			c.updateMessageByIndex(ctx, i)
		case StateTIMEOUT:
			c.updateMessageByIndex(ctx, i)
			c.logger.Printf("task %s timeout\n", message.taskRTName)
		}
	}

}

func (c *Client) restoreInitMessages(ctx context.Context) {
	messages := make([]*Message, 0)
	for i := len(c.messages) - 1; i >= 0; i-- {
		message := c.messages[i]
		if message == nil {
			continue
		}
		if message.taskState == StateINIT {
			messages = append(messages, message)
		}
		c.messages[i] = nil
	}
	if err := c.broker.restore(ctx, messages...); err != nil {
		c.logger.Printf("restore messages err: %v\n", err)
	} else {
		c.logger.Printf("restore %d messages\n", len(messages))
	}
}

func (c *Client) Start(ctx context.Context) error {
	c.setState(ClientRUNNING)
	ticker := time.NewTicker(time.Second * 3)
	stop := make(chan os.Signal)
	signal.Notify(
		stop,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	for {
		select {
		case <-stop:
			signal.Stop(stop)
			c.setState(ClientSTOPPED)
			c.restoreInitMessages(ctx)
		case <-ticker.C:
			c.scanMessages(ctx)
		case <-c.currency.C:
			c.scanMessages(ctx)
		}
	}
}

func (c *Client) Delay(ctx context.Context, tasks []TaskInterface, taskOptions ...TaskOption) error {
	messages := make([]*Message, 0, len(tasks))
	for _, task := range tasks {
		meta := Meta{
			TaskID:   genRandString(8),
			TaskName: getType(task),
		}
		meta.SetDefault()
		for _, taskOption := range taskOptions {
			err := taskOption(&meta)
			if err != nil {
				return err
			}
		}
		message := &Message{
			Meta: meta,
			Task: task,
		}
		if err := c.runTaskHook(ctx, message, BeforeCreate); err != nil {
			return err
		}
		messages = append(messages, message)
	}
	return c.broker.send(ctx, messages...)
}

func (c *Client) runTask(ctx context.Context, message *Message) {
	defer c.currency.release()
	if state := message.taskState; state != StateRUNNING {
		c.logger.Printf("invalid state %d, skip\n", state)
		return
	}
	tsc := make(chan *TaskSignal, 1)
	ctx, cancel := context.WithDeadline(ctx, getNow().Add(message.TaskTimeLimit))
	defer cancel()
	c.runTaskHook(ctx, message, BeforeExecute)
	go func(ctx context.Context, message *Message, tsc chan<- *TaskSignal) {
		defer func() {
			if err := recover(); err != nil {
				tsc <- newTaskSignal(PANIC, err.(error).Error())
			}
		}()
		tsc <- newTaskSignal(START, "")
		if err := message.Task.Execute(ctx); err != nil {
			tsc <- newTaskSignal(ERROR, err.Error())
		} else {
			tsc <- newTaskSignal(DONE, "")
		}
	}(ctx, message, tsc)

	var taskState TaskState
	var ticker = time.NewTicker(time.Millisecond * 100)
label:
	for {
		select {
		case <-ticker.C:
			if c.state == ClientSTOPPED {
				c.runTaskHook(ctx, message, BeforeProcessExit)
				return
			}
		case <-ctx.Done():
			taskState = StateTIMEOUT
			c.runTaskHook(ctx, message, AfterTimeout)
			break label
		case ts := <-tsc:
			// c.logger.Printf("%s %s\n", message.rtName, ts)
			switch ts.Type {
			case DONE:
				taskState = StateSUCCEED
				c.runTaskHook(ctx, message, AfterSucceed)
				break label
			case ERROR:
				c.runTaskHook(ctx, message, AfterFailed)
				taskState = StateFAILED
				break label
			case PANIC:
				c.runTaskHook(ctx, message, AfterFailed)
				taskState = StateFAILED
				break label
			}
		}
	}
	if err := message.SetTaskState(taskState); err != nil {
		c.logger.Printf("set task state err: %v\n", err)
	}
}

func (c *Client) runTaskHook(ctx context.Context, message *Message, taskHook TaskHook) error {
	// 调用方通常无需处理错误
	hooks := message.Task.Hooks(ctx)
	if len(hooks) == 0 {
		return nil
	}
	f := hooks[taskHook]
	if f == nil {
		return nil
	}
	if err := f(ctx); err != nil {
		c.logger.Printf("run %s hook %s err: %v\n", message.taskRTName, taskHook, err)
		return err
	}
	return nil
}
