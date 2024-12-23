package application

import (
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	ucase "github.com/konstantin-kukharev/metrics/domain/usecase/metric"
	handler "github.com/konstantin-kukharev/metrics/internal/controller/rest/metric"
	"github.com/konstantin-kukharev/metrics/internal/repository/memory"
)

type ApplicationConfig interface {
	GetAddress() string
}

type Logger interface {
	Info(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type Server struct {
	config ApplicationConfig
	log    Logger
	server *http.Server
}

func NewServer(app ApplicationConfig, l Logger) *Server {
	store := memory.NewStorage(l)
	add := ucase.NewAddMetric(store)
	getVal := ucase.NewGetMetric(store)
	list := ucase.NewListMetric(store)

	r := chi.NewRouter()
	r.Method("POST", "/update/{type}/{name}/{val}", WithLogging(handler.NewAddMetric(add), l))
	r.Method("GET", "/value/{type}/{name}", WithLogging(handler.NewGetMetric(getVal), l))
	r.Method("GET", "/", WithCompressing(WithLogging(handler.NewIndexMetric(list), l)))

	r.Method("POST", "/update/", WithCompressing(WithLogging(handler.NewAddMetricV2(add), l)))
	r.Method("POST", "/value/", WithCompressing(WithLogging(handler.NewMetricGetV2(getVal), l)))

	return &Server{
		config: app,
		log:    l,
		server: &http.Server{
			Handler: r,
			Addr:    app.GetAddress(),
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func(c context.Context) {
		// blocks here until there's a signal
		<-c.Done()

		err := s.server.Shutdown(c)
		if err != nil {
			s.log.Error("http server shutting down: " + err.Error())
		} else {
			s.log.Info("http server shutdown processed successfully")
		}
	}(ctx)

	return s.server.ListenAndServe()
}

func WithLogging(h http.Handler, l Logger) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)

		duration := time.Since(start)
		l.Info("new request",
			"uri", uri,
			"method", method,
			"duration", duration,
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

// Write is necessary in order to properly implement the io.Writer interface.
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	return n, err
}
