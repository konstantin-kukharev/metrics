package memory

import (
	"context"
	"sync"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type Logger interface {
	Debug(msg string, fields ...any)
	Error(msg string, fields ...any)
}

type key struct {
	t, n string
}

type MetricStorage struct {
	log   Logger
	store map[key]*entity.Metric
	mx    *sync.RWMutex
}

type Updater interface {
	Set(es ...*entity.Metric) error
	Get(*entity.Metric) (*entity.Metric, bool)
}

type CreateOrUpdate func(a Updater) error

func (ms *MetricStorage) SetUnsafe(es ...*entity.Metric) error {
	for _, m := range es {
		k := key{t: m.MType, n: m.ID}
		ms.store[k] = m
	}

	return nil
}

func (ms *MetricStorage) GetUnsafe(m *entity.Metric) (*entity.Metric, bool) {
	k := key{t: m.MType, n: m.ID}
	if v, ok := ms.store[k]; ok {
		return v, ok
	}

	return m, false
}

func (ms *MetricStorage) Set(es ...*entity.Metric) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	return ms.SetUnsafe(es...)
}

func (ms *MetricStorage) Get(m *entity.Metric) (*entity.Metric, bool) {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	return ms.GetUnsafe(m)
}

func (ms *MetricStorage) List() []*entity.Metric {
	list := make([]*entity.Metric, 0, len(ms.store))
	ms.mx.RLock()
	for _, val := range ms.store {
		list = append(list, val)
	}
	ms.mx.RUnlock()

	return list
}

func (ms *MetricStorage) UnitOfWork(ctx context.Context, payload func(context.Context) error) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	return payload(ctx)
}

func NewStorage(l Logger) *MetricStorage {
	ms := new(MetricStorage)
	ms.log = l
	ms.store = map[key]*entity.Metric{}
	ms.mx = &sync.RWMutex{}

	return ms
}