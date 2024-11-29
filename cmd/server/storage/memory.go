package storage

import (
	"sync"

	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/konstantin-kukharev/metrics/pkg/dto"
)

type key struct {
	c, n string
}

type memStorage struct {
	mx      *sync.RWMutex
	storage map[key]string
}

func (s *memStorage) List() []internal.MetricValue {
	res := make([]internal.MetricValue, 0, len(s.storage))
	for k, v := range s.storage {
		res = append(res, dto.NewMetricValue(k.c, k.n, v))
	}

	return res
}

func (s *memStorage) Get(t, k string) (string, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	v, ok := s.storage[key{c: t, n: k}]
	return v, ok
}

func (s *memStorage) Set(t, k string, v string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.storage[key{c: t, n: k}] = v
	return nil
}

func NewMemStorage() *memStorage {
	return &memStorage{
		mx:      &sync.RWMutex{},
		storage: map[key]string{},
	}
}
