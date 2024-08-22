package tinycelery

import "fmt"

type Message struct {
	*Meta
	Task Task `json:"task"`
}

func (m *Message) setDefault() {
	if m.Meta == nil {
		m.Meta = &Meta{}
	}
	m.Meta.setDefault()
	m.Meta.Name = getType(m.Task)
	m.Meta.rtName = fmt.Sprintf("%s[%s]", m.Meta.Name, genRandString(4))
	if m.Meta.RateLimit.Limit != "" && m.Meta.RateLimit.Key == "" {
		m.Meta.RateLimit.Key = m.Meta.Name
	}
}
