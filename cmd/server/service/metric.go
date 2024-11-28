package service

import (
	"github.com/konstantin-kukharev/metrics/internal"
)

type metric struct {
	storage internal.Storage
	metric  map[string]func(s internal.Storage, k, v string) error
}

func (s *metric) Get(k string) ([]byte, bool) {
	return s.storage.Get(k)
}

func (s *metric) Set(t, k string, v string) error {
	if k == "" {
		return internal.ErrWrongMetricName
	}

	if _, ok := s.metric[t]; !ok {
		return internal.ErrWrongMetricType
	}

	return s.metric[t](s.storage, k, v)
}

func NewMetric(s internal.Storage, m ...internal.Metric) *metric {
	srv := &metric{
		storage: s,
		metric:  map[string]func(s internal.Storage, k string, v string) error{},
	}

	for _, ms := range m {
		srv.metric[ms.Name()] = ms.Setter()
	}

	return srv
}
