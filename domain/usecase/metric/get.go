package metric

import (
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type GetMetricProvider interface {
	Get(*entity.Metric) (*entity.Metric, bool)
}

type GetMetric struct {
	provider GetMetricProvider
}

func (am *GetMetric) Do(m *entity.Metric) (*entity.Metric, bool) {
	return am.provider.Get(m)
}

func NewGetMetric(s GetMetricProvider) *GetMetric {
	srv := &GetMetric{
		provider: s,
	}

	return srv
}
