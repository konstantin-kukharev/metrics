package metric

import (
	"context"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type Tx interface {
	UnitOfWork(context.Context, func(context.Context) error) error
}

type UnsafeAddMetricProvider interface {
	SetUnsafe(m ...*entity.Metric) error
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
			res, ok := am.provider.GetUnsafe(m)
			if ok {
				m.Aggregate(res)
			}

			if err := am.provider.SetUnsafe(m); err != nil {
				return err
			}
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
