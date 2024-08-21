package tinycelery

type initOption func(c *Client) error

func WithCurrencyMax(currencyMax uint32) initOption {
	return func(c *Client) error {
		if currencyMax < 1 {
			return ErrInvalidOption
		}
		c.currency.max = currencyMax
		return nil
	}
}

func WithPrefetch(prefetch uint32) initOption {
	return func(c *Client) error {
		c.prefetch = prefetch
		return nil
	}
}

func WithQueue(queue string) initOption {
	return func(c *Client) error {
		c.broker.queue = queue
		return nil
	}
}
