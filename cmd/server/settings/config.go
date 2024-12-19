package settings

import (
	"flag"
	"os"

	"github.com/konstantin-kukharev/metrics/internal"
)

type Config struct {
	Address string
}

func (c *Config) GetAddress() string {
	return c.Address
}

func NewConfig() *Config {
	c := &Config{Address: internal.DefaultServerAddr}

	return c
}

func (c *Config) WithFlag() *Config {
	flag.StringVar(&c.Address, "a", internal.DefaultServerAddr, "server address")
	flag.Parse()

	return c
}

func (c *Config) WithEnv() *Config {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		c.Address = envRunAddr
	}

	return c
}
