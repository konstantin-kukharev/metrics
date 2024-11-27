package storage

import "sync"

type memStorage struct {
	mx      *sync.Mutex
	storage map[string][]byte
}

func (s *memStorage) Get(k string) ([]byte, bool) {
	v, ok := s.storage[k]
	return v, ok
}

func (s *memStorage) Set(k string, v []byte) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.storage[k] = v
}

func NewMemStorage() *memStorage {
	return &memStorage{
		mx:      &sync.Mutex{},
		storage: map[string][]byte{},
	}
}
