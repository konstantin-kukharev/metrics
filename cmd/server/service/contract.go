package service

import "github.com/konstantin-kukharev/metrics/pkg/metric"

type Metric interface {
	List() []metric.Value
	Get(t, k string) ([]byte, bool)
	Set(t, k string, v string) error
}
