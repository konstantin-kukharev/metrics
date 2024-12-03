package storage

import "github.com/konstantin-kukharev/metrics/pkg/metric"

type Store interface {
	List() []metric.Value
	Get(t metric.Type, k string) (string, bool)
	Set(t metric.Type, k string, v string) error
}
