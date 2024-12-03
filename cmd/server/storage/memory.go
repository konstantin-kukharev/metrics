package storage

import (
	"sync"

	"github.com/konstantin-kukharev/metrics/pkg/metric"
)

type key struct {
	c metric.Type
	n string
}

type memStorage struct {
	mx      *sync.RWMutex
	storage map[key]string
}

func (s *memStorage) List() []metric.Value {
	res := make([]metric.Value, 0, len(s.storage))
	for k, v := range s.storage {
		val, err := metric.NewValue(k.c, k.n, v)
		//todo: catch err
		if err != nil {
			continue
		}
		res = append(res, val)
	}

	return res
}

func (s *memStorage) Get(t metric.Type, k string) (string, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	v, ok := s.storage[key{c: t, n: k}]
	return v, ok
}

func (s *memStorage) Set(t metric.Type, k string, v string) error {
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
