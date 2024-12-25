package event

import "github.com/konstantin-kukharev/metrics/domain/entity"

type MetricAdd struct {
	Metric *entity.Metric
}
