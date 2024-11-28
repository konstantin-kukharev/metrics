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

type Handler interface {
	MetricUpdate(w http.ResponseWriter, r *http.Request)
}

type Metric interface {
	Name() string
	Setter() func(s Storage, k, v string) error
}
