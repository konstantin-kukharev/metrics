package roundtripper

import (
	"net/http"
	"time"

	"github.com/konstantin-kukharev/metrics/internal/logger"
	"go.uber.org/zap"
)

type Logging struct {
	next http.RoundTripper
	log  *logger.ZapLogger
}

func NewLogging(next http.RoundTripper, l *logger.ZapLogger) *Logging {
	return &Logging{
		next: next,
		log:  l,
	}
}

func (rt *Logging) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	defer func(begin time.Time) {
		rt.log.InfoCtx(req.Context(), "Request",
			zap.String("method", req.Method),
			zap.String("host", req.URL.Scheme+"://"+req.URL.Host+req.URL.Path),
			zap.Any("error", err),
			zap.Duration("took", time.Since(begin)),
		)
	}(time.Now())

	return rt.next.RoundTrip(req)
}
