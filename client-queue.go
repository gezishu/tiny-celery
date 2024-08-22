package tinycelery

func WithQueue(queue string) clientInitOption {
	return func(c *Client) {
		c.broker.queue = queue
	}
}
