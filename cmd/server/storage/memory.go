package storage

import (
	"sync"

	"github.com/konstantin-kukharev/metrics/internal/metric"
)

type Store interface {
	List() []metric.Value
	Get(t metric.Type, k string) (string, bool)
	Set(t metric.Type, k string, v string) error
}

type key struct {
	c metric.Type
	n string
}

type MemStorage struct {
	mx      *sync.RWMutex
	storage map[key]string
}

func (s *MemStorage) List() []metric.Value {
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

func (s *MemStorage) Get(t metric.Type, k string) (string, bool) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	v, ok := s.storage[key{c: t, n: k}]
	return v, ok
}

func (s *MemStorage) Set(t metric.Type, k string, v string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.storage[key{c: t, n: k}] = v
	return nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		mx:      &sync.RWMutex{},
		storage: map[key]string{},
	}
}
