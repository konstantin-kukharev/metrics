package settings

import (
	"flag"

	"github.com/konstantin-kukharev/metrics/internal"
)

func fromFlag(s *Config) {
	s.Address = *flag.String("a", internal.DefaultServerAddr, "server address")

	flag.Parse()
}
