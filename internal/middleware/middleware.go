package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/konstantin-kukharev/metrics/internal/logger"
	"go.uber.org/zap"
)

type LoggingRoundTripper struct {
	next http.RoundTripper
	log  *logger.ZapLogger
}

type CompressRoundTripper struct {
	next http.RoundTripper
}

func NewLoggingRoundTripper(next http.RoundTripper, l *logger.ZapLogger) *LoggingRoundTripper {
	return &LoggingRoundTripper{
		next: next,
		log:  l,
	}
}

func (rt *LoggingRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
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

func NewCompressRoundTripper(next http.RoundTripper) *CompressRoundTripper {
	return &CompressRoundTripper{
		next: next,
	}
}

func (rt *CompressRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	var buf bytes.Buffer
	g := gzip.NewWriter(&buf)
	b, err := io.ReadAll(req.Body)
	defer req.Body.Close()

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
	r, err := http.NewRequestWithContext(req.Context(), req.Method, url, &buf)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Accept-Encoding", "gzip")
	r.Header.Set("Content-Type", req.Header.Get("Content-Type"))

	return rt.next.RoundTrip(r)
}

func WithLogging(h http.Handler, l *logger.ZapLogger) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)

		duration := time.Since(start)
		l.InfoCtx(r.Context(), "new request",
			zap.String("uri", uri),
			zap.String("method", method),
			zap.Duration("duration", duration),
		)
	}

	return http.HandlerFunc(logFn)
}

func WithCompressing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := r.Header.Get("Content-Encoding")
		if encoding == "gzip" {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer reader.Close()

			r.Body = io.NopCloser(reader)
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			h.ServeHTTP(w, r)
			return
		}
		defer gz.Close()

		gzrw := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		h.ServeHTTP(gzrw, r)
	})
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	return n, err
}
