package settings

import (
	"os"
)

// ADDRESS отвечает за адрес эндпоинта HTTP-сервера.
func fromEnv(s *Config) {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		s.Address = envRunAddr
	}
}
