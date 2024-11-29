package service

import (
	"github.com/konstantin-kukharev/metrics/internal"
)

type metric struct {
	storage  internal.Storage
	metric   map[string]internal.Metric
	cacheKey string
	cache    []internal.MetricValue
}

func (s *metric) List() []internal.MetricValue {
	val, ok := s.storage.Get(internal.MetricCounter, internal.CacheKey)
	if ok && s.cacheKey != val {
		s.cache = s.storage.List()
	}

	return s.cache
}

func (s *metric) Get(t, k string) ([]byte, bool) {
	if _, ok := s.metric[t]; !ok {
		return []byte{}, ok
	}
	val, ok := s.storage.Get(t, k)
	if !ok {
		return []byte{}, ok
	}

	return []byte(val), true
}

func (s *metric) Set(t, k string, v string) error {
	var resultVal string
	if k == "" {
		return internal.ErrWrongMetricName
	}

	if _, ok := s.metric[t]; !ok {
		return internal.ErrWrongMetricType
	}

	val, err := s.metric[t].Encode(v)
	if err != nil {
		return internal.ErrInvalidData
	}

	var val1 []byte
	vs, ok := s.storage.Get(t, k)
	if !ok {
		resultVal, err = s.metric[t].Decode(val)
		if err != nil {
			return internal.ErrInvalidData
		}

		return s.storage.Set(t, k, resultVal)
	}

	if val1, err = s.metric[t].Encode(vs); err != nil {
		return internal.ErrInvalidData
	}

	val, err = s.metric[t].Addition(val, val1)

	if err != nil {
		return internal.ErrInvalidData
	}

	resultVal, err = s.metric[t].Decode(val)

	if err != nil {
		return internal.ErrInvalidData
	}

	return s.storage.Set(t, k, resultVal)
}

func NewMetric(s internal.Storage, m ...internal.Metric) *metric {
	srv := &metric{
		storage: s,
		metric:  map[string]internal.Metric{},
	}

	for _, ms := range m {
		srv.metric[ms.Name()] = ms
	}

	return srv
}
