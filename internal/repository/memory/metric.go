package memory

import (
	"context"
	"encoding/json"
	"os"
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
	log    Logger
	store  map[key]*entity.Metric
	mx     *sync.RWMutex
	backup *os.File
}

func (ms *MetricStorage) GetUnsafe(m *entity.Metric) (*entity.Metric, bool) {
	k := key{t: m.MType, n: m.ID}
	if v, ok := ms.store[k]; ok {
		return v, ok
	}

	return m, false
}

func (ms *MetricStorage) ListUnsafe() []*entity.Metric {
	list := make([]*entity.Metric, 0, len(ms.store))
	for _, val := range ms.store {
		list = append(list, val)
	}

	return list
}

func (ms *MetricStorage) Set(es ...*entity.Metric) {
	for _, m := range es {
		k := key{t: m.MType, n: m.ID}
		ms.store[k] = m
		if ms.backup != nil {
			if b, err := json.Marshal(m); err == nil {
				_, _ = ms.backup.Write(b)
				_, _ = ms.backup.WriteString("\n")
			}
		}
	}
}

func (ms *MetricStorage) Get(m *entity.Metric) (*entity.Metric, bool) {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	return ms.GetUnsafe(m)
}

func (ms *MetricStorage) List() []*entity.Metric {
	ms.mx.RLock()
	defer ms.mx.RUnlock()

	return ms.ListUnsafe()
}

func (ms *MetricStorage) UnitOfWork(ctx context.Context, payload func(context.Context) error) error {
	ms.mx.Lock()
	defer ms.mx.Unlock()
	return payload(ctx)
}

func (ms *MetricStorage) WithStream(f *os.File) {
	ms.backup = f
}

func NewStorage(l Logger) *MetricStorage {
	ms := new(MetricStorage)
	ms.log = l
	ms.store = map[key]*entity.Metric{}
	ms.mx = &sync.RWMutex{}

	return ms
}
