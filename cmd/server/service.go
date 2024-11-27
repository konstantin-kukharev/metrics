package main

import (
	"github.com/konstantin-kukharev/metrics/cmd/server/internal"
	"github.com/konstantin-kukharev/metrics/cmd/server/metric"
)

type service struct {
	storage internal.Storage
	metric  map[string]func(s internal.Storage, k, v string) error
}

func (s *service) Get(k string) ([]byte, bool) {
	return s.storage.Get(k)
}

func (s *service) Set(t, k string, v string) error {
	if k == "" {
		return internal.ErrWrongMetricName
	}

	if _, ok := s.metric[t]; !ok {
		return internal.ErrWrongMetricType
	}

	return s.metric[t](s.storage, k, v)
}

func NewMetricService(s internal.Storage, m ...metric.Metric) *service {
	srv := &service{
		storage: s,
		metric:  map[string]func(s internal.Storage, k string, v string) error{},
	}

	for _, ms := range m {
		srv.metric[ms.Name] = ms.Setter
	}

	return srv
}
