package settings

import (
	"os"
	"strconv"
)

// ADDRESS отвечает за адрес эндпоинта HTTP-сервера.
// REPORT_INTERVAL позволяет переопределять reportInterval.
// POLL_INTERVAL позволяет переопределять pollInterval.
func fromEnv(s *Config) {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		s.Address = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if val, err := strconv.Atoi(envReportInterval); err == nil {
			s.ReportInterval = val
		}
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval != "" {
		if val, err := strconv.Atoi(envPoolInterval); err == nil {
			s.PoolInterval = val
		}
	}
}
