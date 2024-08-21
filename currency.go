package tinycelery

import "sync"

type Currency struct {
	sync.Mutex
	max  uint32
	curr uint32
	C    chan uint32 // only at release
}

func (c *Currency) prerequire() bool {
	return c.curr < c.max
}

func (c *Currency) require() bool {
	c.Lock()
	defer c.Unlock()
	if c.curr < c.max {
		c.curr++
		return true
	}
	return false
}

func (c *Currency) release() {
	c.Lock()
	defer c.Unlock()
	c.curr--
	go func() {
		c.C <- c.curr
	}()
}
