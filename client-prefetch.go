package tinycelery

const defaultPrefetch uint32 = 0

func WithPrefetch(prefetch uint32) clientInitOption {
	return func(c *Client) {
		c.prefetch = prefetch
	}
}
