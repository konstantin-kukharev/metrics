package metric

import (
	"context"
	"errors"

	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type ReportProvider interface {
	Do(*entity.Metric) error
}

type ReportMetricProvider interface {
	Tx
	ListUnsafe() []*entity.Metric
}

type ReportMetric struct {
	provider ReportMetricProvider
	reporter ReportProvider
}

func (am *ReportMetric) Do() error {
	return am.provider.UnitOfWork(context.TODO(), func(_ context.Context) error {
		var errs error
		for _, m := range am.provider.ListUnsafe() {
			err := am.reporter.Do(m)
			errs = errors.Join(errs, err)
		}

		return errs
	})
}

func NewReportMetric(a ReportMetricProvider, r ReportProvider) *ReportMetric {
	srv := &ReportMetric{
		provider: a,
		reporter: r,
	}

	return srv
}
