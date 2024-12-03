package settings

import (
	"flag"

	"github.com/konstantin-kukharev/metrics/internal"
)

func fromFlag(s *Config) {
	flag.StringVar(&s.Address, "a", internal.DefaultServerAddr, "server address")
	flag.DurationVar(&s.ReportInterval, "r", internal.DefaultReportInterval, "report interval time duration")
	flag.DurationVar(&s.PoolInterval, "p", internal.DefaultPoolInterval, "pool interval time duration")

	flag.Parse()
}
