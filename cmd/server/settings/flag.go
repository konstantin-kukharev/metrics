package settings

import (
	"flag"

	"github.com/konstantin-kukharev/metrics/internal"
)

func fromFlag(s *settings) {
	s.address = *flag.String("a", internal.DefaultServerAddr, "server address")

	flag.Parse()
}
