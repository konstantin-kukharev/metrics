package metric

import (
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type ListMetricProvider interface {
	List() []*entity.Metric
}

type ListMetric struct {
	provider ListMetricProvider
}

func (am *ListMetric) Do() []*entity.Metric {
	return am.provider.List()
}

func NewListMetric(s ListMetricProvider) *ListMetric {
	srv := &ListMetric{
		provider: s,
	}

	return srv
}
