package settings

import (
	"flag"

	"github.com/konstantin-kukharev/metrics/internal"
)

func fromFlag(s *settings) {
	s.address = *flag.String("a", internal.DefaultServerAddr, "server address")
	s.reportInterval = *flag.Duration("r", internal.DefaultReportInterval, "report interval time duration in seconds")
	s.poolInterval = *flag.Duration("p", internal.DefaultPoolInterval, "pool interval time duration in seconds")

	flag.Parse()
}
