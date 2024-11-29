package internal

import (
	"runtime"
)

type Storage interface {
	List() []MetricValue
	Get(t, k string) (string, bool)
	Set(t, k string, v string) error
}

type MetricService interface {
	List() []MetricValue
	Get(t, k string) ([]byte, bool)
	Set(t, k string, v string) error
}

type Metric interface {
	Name() string
	Encode(v string) ([]byte, error)
	Decode(v []byte) (string, error)
	//Addition FiFo or summ
	Addition(...[]byte) ([]byte, error)
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
