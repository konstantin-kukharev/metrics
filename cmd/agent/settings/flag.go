package settings

import (
	"flag"
	"time"

	"github.com/konstantin-kukharev/metrics/internal"
)

func fromFlag(s *Config) {
	var pi, ri int
	flag.StringVar(&s.Address, "a", internal.DefaultServerAddr, "server address")
	flag.IntVar(&ri, "r", 10, "report interval time duration")
	flag.IntVar(&pi, "p", 2, "pool interval time duration")

	s.ReportInterval = time.Duration(ri)
	s.PoolInterval = time.Duration(pi)

	flag.Parse()
}
