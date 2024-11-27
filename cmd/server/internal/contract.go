package internal

import "net/http"

type Storage interface {
	Get(k string) ([]byte, bool)
	Set(k string, v []byte)
}

type MetricService interface {
	Get(k string) ([]byte, bool)
	Set(t, k string, v string) error
}

type Server interface {
	MetricUpdate(w http.ResponseWriter, r *http.Request)
}
