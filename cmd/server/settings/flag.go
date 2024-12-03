package settings

import (
	"flag"

	"github.com/konstantin-kukharev/metrics/internal"
)

func fromFlag(s *Config) {
	flag.StringVar(&s.Address, "a", internal.DefaultServerAddr, "server address")

	flag.Parse()
}
