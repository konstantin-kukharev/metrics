package service

import (
	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/konstantin-kukharev/metrics/pkg/metric"
)

type Metric interface {
	List() []metric.Value
	Get(t, k string) ([]byte, bool)
	Set(t, k string, v string) error
}

type MetricService struct {
	cfg      settings.Settings
	storage  storage.Store
	metric   map[string]metric.Type
	cacheKey string
	cache    []metric.Value
}

func (s *MetricService) List() []metric.Value {
	val, ok := s.storage.Get(&metric.Counter{}, internal.CacheKey)
	if ok && s.cacheKey != val {
		s.cache = s.storage.List()
	}

	return s.cache
}

func (s *MetricService) Get(t, k string) ([]byte, bool) {
	if _, ok := s.metric[t]; !ok {
		return []byte{}, ok
	}
	val, ok := s.storage.Get(s.metric[t], k)
	if !ok {
		return []byte{}, ok
	}

	return []byte(val), true
}

func (s *MetricService) Set(t, k string, v string) error {
	//валидация параметров
	if k == "" {
		return internal.ErrWrongMetricName
	}

	if _, ok := s.metric[t]; t == "" || !ok {
		return internal.ErrWrongMetricType
	}

	val, err := s.metric[t].Encode(v)
	if err != nil {
		return internal.ErrInvalidData
	}

	valStore, ok := s.storage.Get(s.metric[t], k)
	if !ok {
		resp, err := s.metric[t].Decode(val)
		if err != nil {
			return internal.ErrInvalidData
		}

		return s.storage.Set(s.metric[t], k, resp)
	}

	val1, err := s.metric[t].Encode(valStore)
	if err != nil {
		return internal.ErrInvalidData
	}

	val2, err := s.metric[t].Aggregate(val, val1)
	if err != nil {
		return internal.ErrInvalidData
	}

	resultVal, err := s.metric[t].Decode(val2)
	if err != nil {
		return internal.ErrInvalidData
	}

	return s.storage.Set(s.metric[t], k, resultVal)
}

func NewMetric(cfg settings.Settings, s storage.Store, m ...metric.Type) *MetricService {
	srv := &MetricService{
		cfg:     cfg,
		storage: s,
		metric:  map[string]metric.Type{},
	}

	for _, ms := range m {
		srv.metric[ms.GetName()] = ms
	}

	return srv
}
