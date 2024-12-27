package roundtripper

import (
	"net/http"
	"time"
)

type logger interface {
	Info(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type Logging struct {
	next http.RoundTripper
	log  logger
}

func NewLogging(next http.RoundTripper, l logger) *Logging {
	return &Logging{
		next: next,
		log:  l,
	}
}

func (rt *Logging) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	defer func(begin time.Time) {
		rt.log.Info("Request",
			"method", req.Method,
			"host", req.URL.Scheme+"://"+req.URL.Host+req.URL.Path,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	return rt.next.RoundTrip(req)
}
