package settings

import (
	"flag"

	"github.com/konstantin-kukharev/metrics/internal"
)

func fromFlag(s *Config) {
	s.Address = *flag.String("a", internal.DefaultServerAddr, "server address")
	s.ReportInterval = *flag.Duration("r", internal.DefaultReportInterval, "report interval time duration in seconds")
	s.PoolInterval = *flag.Duration("p", internal.DefaultPoolInterval, "pool interval time duration in seconds")

	flag.Parse()
}
