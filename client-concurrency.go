package tinycelery

import (
	"sync"
	"time"
)

func WithConcurrency(concurrency uint32) clientInitOption {
	return func(c *Client) {
		if concurrency == 0 {
			panic("invalid concurrency")
		}
		c.concurrency.max = concurrency
	}
}

var defaultConcurrency = &concurrency{
	max: 1,
	cur: 0,
	c:   make(chan int, 10),
}

type concurrency struct {
	sync.Mutex
	max uint32
	cur uint32
	// cur每次的变化会被写入c
	c chan int
}

func (c *concurrency) writeC(i int) {
	select {
	case c.c <- i:
		return
	case <-time.After(time.Millisecond):
		// 主要是用于快速触发scanMessages, 允许丢失
		return
	}
}

func (c *concurrency) require(pre bool) bool {
	if pre {
		return c.cur < c.max
	}
	c.Lock()
	defer c.Unlock()
	if c.cur < c.max {
		c.cur++
		c.writeC(1)
		return true
	}
	return false
}

func (c *concurrency) release() {
	c.Lock()
	defer c.Unlock()
	if c.cur > 0 {
		c.cur--
		c.writeC(-1)
	}
}
