package memory

import (
	"sync"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type Logger interface {
	Debug(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type Write interface {
	Set(m *entity.Metric) error
}

type Read interface {
	Get(*entity.Metric) (*entity.Metric, bool)
}

type List interface {
	List() []*entity.Metric
}

type AddMetricProvider interface {
	Write
	Read
	List

	CreateOrUpdate(func(a Write, g Read) error) error
}

type key struct {
	t, n string
}

type MetricStorage struct {
	log   Logger
	store map[key]*entity.Metric
	sync.RWMutex
}

func (ms *MetricStorage) Set(m *entity.Metric) error {
	k := key{t: m.MType, n: m.ID}
	ms.Lock()
	ms.store[k] = m
	ms.Unlock()

	return nil
}

func (ms *MetricStorage) Get(m *entity.Metric) (*entity.Metric, bool) {
	k := key{t: m.MType, n: m.ID}
	ms.RLock()
	if v, ok := ms.store[k]; ok {
		return v, ok
	}
	ms.RUnlock()

	return m, false
}

func (ms *MetricStorage) List() []*entity.Metric {
	list := make([]*entity.Metric, 0, len(ms.store))
	ms.RLock()
	for _, val := range ms.store {
		list = append(list, val)
	}
	ms.RUnlock()

	return list
}

func (ms *MetricStorage) CreateOrUpdate(payload func(a Write, g Read) error) error {
	ms.Lock()
	defer ms.Unlock()

	return payload(ms, ms)
}

func NewStorage(l Logger) *MetricStorage {
	ms := new(MetricStorage)
	ms.log = l
	ms.store = map[key]*entity.Metric{}

	return ms
}
