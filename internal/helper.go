package internal

import (
	"strconv"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

func GetMetricValue(m *entity.Metric) string {
	switch m.MType {
	case domain.MetricGauge:
		return strconv.FormatFloat(*m.Value, 'f', DefaultFloatPrecision, 64)
	case domain.MetricCounter:
		return strconv.FormatInt(*m.Delta, 10)
	}

	return ""
}
