package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"time"
)

type LoggingRoundTripper struct {
	next http.RoundTripper
	log  Logger
}

type CompressRoundTripper struct {
	next http.RoundTripper
}

func NewLoggingRoundTripper(next http.RoundTripper, l Logger) *LoggingRoundTripper {
	return &LoggingRoundTripper{
		next: next,
		log:  l,
	}
}

func (rt *LoggingRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
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

func NewCompressRoundTripper(next http.RoundTripper) *CompressRoundTripper {
	return &CompressRoundTripper{
		next: next,
	}
}

func (rt *CompressRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	if _, err = g.Write(b); err != nil {
		return
	}
	if err = g.Close(); err != nil {
		return
	}

	url := req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
	r, err := http.NewRequest(req.Method, url, &buf)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("Content-Type", req.Header.Get("Content-Type"))

	return rt.next.RoundTrip(r)
}
