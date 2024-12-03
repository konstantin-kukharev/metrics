package settings

import (
	"flag"

	"github.com/konstantin-kukharev/metrics/internal"
)

func fromFlag(s *Config) {
	flag.StringVar(&s.Address, "a", internal.DefaultServerAddr, "server address")
	flag.IntVar(&s.ReportInterval, "r", 10, "report interval time duration")
	flag.IntVar(&s.PoolInterval, "p", 2, "pool interval time duration")

	flag.Parse()
}
