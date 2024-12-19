package entity

import (
	"strconv"

	"github.com/konstantin-kukharev/metrics/domain"
)

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metric) Aggregate(m2 *Metric) {
	*m.Delta = *m2.Delta + *m.Delta
	*m.Value = *m2.Value
}

func NewMetric(name, mtype, value string) (*Metric, error) {
	m := new(Metric)
	if name == "" {
		return m, domain.ErrWrongMetricName
	}

	switch m.MType {
	case domain.MetricGauge:
		cv, err := strconv.ParseFloat(value, 64)
		if err == nil {
			*m.Value = cv
			return m, nil
		}
		return m, domain.ErrInvalidData
	case domain.MetricCounter:
		iv, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			*m.Delta = iv
			return m, nil
		}
		return m, domain.ErrInvalidData
	}

	return m, domain.ErrWrongMetricType
}
