package internal

import (
	"net/http"
	"runtime"
)

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

type MetricValue interface {
	Type() string
	Name() string
	Value() string
}

type StateMemory interface {
	Update(m *runtime.MemStats)
	Data() []MetricValue
}

type AgentReporter interface {
	Report([]MetricValue)
}
