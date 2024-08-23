package tinycelery

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

type clientInitOption func(c *Client)

type Client struct {
	sync.Mutex
	concurrency *concurrency
	prefetch    uint32
	broker      *redisBroker
	logger      *log.Logger
	state       clientState
	messages    []*Message
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetAccumulate(ctx context.Context) (uint32, error) {
	return c.broker.getAccumulate(ctx)
}

func (c *Client) setState(state clientState) {
	c.Lock()
	defer c.Unlock()
	c.state = state
}

func (c *Client) Init(rc *redis.Client, queue string, options ...clientInitOption) (*Client, error) {
	c.concurrency = newDefaultConcurrency()
	c.prefetch = defaultPrefetch
	c.logger = log.New(os.Stdout, "[tiny-celery] ", log.Ltime)
	c.broker = newDefaultRedisBroker()
	c.broker.rc = rc
	c.broker.queue = queue
	for _, option := range options {
		option(c)
	}
	c.messages = make([]*Message, c.concurrency.max+c.prefetch)
	return c, nil
}

func (c *Client) RegisterTask(task Task) error {
	taskName := getType(task)
	if _, registered := c.broker.tasks[taskName]; registered {
		return ErrTaskHasRegistered
	}
	c.broker.tasks[taskName] = task
	return nil
}

func (c *Client) updateMessage(ctx context.Context, index int) error {
	c.messages[index] = nil
	message, err := c.broker.fetch(ctx)
	if err != nil {
		return err
	}
	c.messages[index] = message
	return nil
}

func (c *Client) scan(ctx context.Context) {
	if c.state != clientRUNNING {
		return
	}
	for i := 0; i < len(c.messages); i++ {
		message := c.messages[i]
		if message == nil {
			continue
		}
		switch message.Meta.state {
		case taskINIT:
			if message.Meta.ETA > getNow().Unix() {
				continue
			}
			if !c.concurrency.require(true) {
				break
			}
			if message.Meta.RateLimit.limited(ctx, c.broker.rc) {
				continue
			}
			if c.concurrency.require(false) {
				if err := message.Meta.setState(taskRUNNING); err != nil {
					panic(fmt.Errorf("unreach err: %v", err))
				}
				go c.runTask(ctx, message)
			}
		case taskSUCCEED:
			_ = c.updateMessage(ctx, i)
		case taskFAILED:
			// todo 增加重试逻辑
			c.logger.Printf("task %s failed\n", message.rtName)
			_ = c.updateMessage(ctx, i)
		case taskTIMEOUT:
			_ = c.updateMessage(ctx, i)
			c.logger.Printf("task %s timeout\n", message.rtName)
		}
	}
	for i := 0; i < len(c.messages); i++ {
		if c.messages[i] == nil {
			if err := c.updateMessage(ctx, i); err != nil {
				if errors.Is(err, ErrMessageNotExists) {
					break
				} else {
					c.logger.Printf("update message err: %v\n", err)
				}
				continue
			}
		}
	}

}

func (c *Client) restore(ctx context.Context) {
	messages := make([]*Message, 0)
	for i := len(c.messages) - 1; i >= 0; i-- {
		message := c.messages[i]
		if message == nil {
			continue
		}
		if message.state != taskINIT {
			continue
		}
		messages = append(messages, message)
		c.messages[i] = nil
	}
	if err := c.broker.send(ctx, messages...); err != nil {
		c.logger.Printf("restore messages err: %v\n", err)
	} else {
		c.logger.Printf("restore %d messages\n", len(messages))
	}
}

func (c *Client) Start(ctx context.Context) error {
	c.setState(clientRUNNING)
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT) //nolint:all
	for {
		select {
		case <-stop:
			signal.Stop(stop)
			c.setState(clientSTOPPED)
			c.restore(ctx)
			return nil
		case <-ticker.C:
			c.scan(ctx)
		case n := <-c.concurrency.c:
			if n < 0 {
				c.scan(ctx)
			}
		}
	}
}

func (c *Client) Delay(ctx context.Context, tasks []Task, options ...taskOption) error {
	messages := make([]*Message, 0, len(tasks))
	for _, task := range tasks {
		message := &Message{
			Task: task,
		}
		message.setDefault()
		for _, option := range options {
			err := option(message.Meta)
			if err != nil {
				return err
			}
		}
		if err := c.runTaskHook(ctx, message, BeforeCreate); err != nil {
			return err
		}
		messages = append(messages, message)
	}
	return c.broker.send(ctx, messages...)
}

func (c *Client) runTask(ctx context.Context, message *Message) {
	defer c.concurrency.release()
	if state := message.state; state != taskRUNNING {
		c.logger.Printf("invalid state %d, skip\n", state)
		return
	}
	tsc := make(chan *TaskSignal, 1)
	ctx, cancel := context.WithDeadline(ctx, getNow().Add(message.TimeLimit))
	defer cancel()
	_ = c.runTaskHook(ctx, message, BeforeExecute)
	go func(ctx context.Context, message *Message, tsc chan<- *TaskSignal) {
		defer func() {
			if err := recover(); err != nil {
				tsc <- newTaskSignal(PANIC, fmt.Sprintf("%v", err))
			}
		}()
		tsc <- newTaskSignal(START, "")
		if err := message.Task.Execute(ctx); err != nil {
			tsc <- newTaskSignal(ERROR, err.Error())
		} else {
			tsc <- newTaskSignal(DONE, "")
		}
	}(ctx, message, tsc)

	var state taskState
	var ticker = time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
label:
	for {
		select {
		case <-ticker.C:
			if c.state == clientSTOPPED {
				_ = c.runTaskHook(ctx, message, BeforeProcessExit)
				return
			}
		case <-ctx.Done():
			state = taskTIMEOUT
			_ = c.runTaskHook(ctx, message, AfterTimeout)
			break label
		case ts := <-tsc:
			switch ts.Type {
			case DONE:
				state = taskSUCCEED
				_ = c.runTaskHook(ctx, message, AfterSucceed)
				break label
			case ERROR, PANIC:
				_ = c.runTaskHook(ctx, message, AfterFailed)
				state = taskFAILED
				break label
			}
		}
	}
	err := message.Meta.setState(state)
	if err != nil {
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
		c.logger.Printf("run %s hook %s err: %v\n", message.rtName, taskHook, err)
		return err
	}
	return nil
}
