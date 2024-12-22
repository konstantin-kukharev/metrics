package metric

import (
	"context"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type Tx interface {
	UnitOfWork(context.Context, func(context.Context) error) error
}

type AddMetricProvider interface {
	Set(m ...*entity.Metric) error
}

type UpdateMetricProvider interface {
	Tx
	AddMetricProvider
	GetMetricProvider
}

type AddMetric struct {
	provider UpdateMetricProvider
}

func (am *AddMetric) Do(ms ...*entity.Metric) error {
	return am.provider.UnitOfWork(context.TODO(), func(_ context.Context) error {
		for _, m := range ms {
			res, ok := am.provider.Get(m)
			if ok {
				res.Aggregate(m)
			} else {
				res = m
			}

			if err := am.provider.Set(res); err != nil {
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
