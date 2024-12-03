package settings

import "github.com/konstantin-kukharev/metrics/internal"

type Settings interface {
	GetAddress() string
}

type Config struct {
	Address string
}

func (c *Config) GetAddress() string {
	return c.Address
}

func New() *Config {
	s := &Config{Address: internal.DefaultServerAddr}
	fromFlag(s)

	return s
}
