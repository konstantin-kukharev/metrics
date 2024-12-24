package metric

import (
	"context"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type Tx interface {
	UnitOfWork(context.Context, func(context.Context) error) error
}

type UnsafeAddMetricProvider interface {
	Set(m ...*entity.Metric)
}

type UnsafeGetMetricProvider interface {
	GetUnsafe(*entity.Metric) (*entity.Metric, bool)
}

type UpdateMetricProvider interface {
	Tx
	UnsafeAddMetricProvider
	UnsafeGetMetricProvider
}

type AddMetric struct {
	provider UpdateMetricProvider
}

func (am *AddMetric) Do(ms ...*entity.Metric) error {
	return am.provider.UnitOfWork(context.TODO(), func(_ context.Context) error {
		for _, m := range ms {
			if res, ok := am.provider.GetUnsafe(m); ok {
				m.Aggregate(res)
			}

			am.provider.Set(m)
		}

		return nil
	})
}

func NewAddMetric(a UpdateMetricProvider) *AddMetric {
	srv := &AddMetric{
		provider: a,
	}

	return srv
}
