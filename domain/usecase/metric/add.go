package metric

import (
	"errors"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type Tx interface {
	CreateOrUpdate(func(a AddMetricProvider, g GetMetricProvider) error) error
}

type AddMetricProvider interface {
	Set(m *entity.Metric) error
	Tx
}

type AddMetric struct {
	addProvider AddMetricProvider
}

func (am *AddMetric) Do(ms ...*entity.Metric) error {
	return am.addProvider.CreateOrUpdate(func(a AddMetricProvider, g GetMetricProvider) error {
		var errs error
		for _, m := range ms {
			if res, ok := g.Get(m); ok {
				res.Aggregate(m)
			}

			if err := a.Set(m); err != nil {
				errors.Join(errs, err)
			}
		}
		return errs
	})
}

func NewAddMetric(a AddMetricProvider) *AddMetric {
	srv := &AddMetric{
		addProvider: a,
	}

	return srv
}
