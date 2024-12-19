package metric

import (
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type ListMetricProvider interface {
	List() ([]*entity.Metric, bool)
}

type ListMetric struct {
	provider ListMetricProvider
}

func (am *ListMetric) Do() ([]*entity.Metric, bool) {
	return am.provider.List()
}

func NewListMetric(s ListMetricProvider) *ListMetric {
	srv := &ListMetric{
		provider: s,
	}

	return srv
}
