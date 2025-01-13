package entity

import (
	"strconv"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/internal"
)

type Metric struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
	MValue
}

type MValue struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metric) Aggregate(m2 *Metric) {
	if m.Delta != nil && m2.Delta != nil {
		*m.Delta += *m2.Delta
	}
}

func (m *Metric) GetValue() string {
	switch m.MType {
	case domain.MetricGauge:
		if m.Value == nil {
			break
		}

		return strconv.FormatFloat(*m.Value, 'f', internal.DefaultFloatPrecision, 64)
	case domain.MetricCounter:
		if m.Delta == nil {
			break
		}

		return strconv.FormatInt(*m.Delta, 10)
	}

	return ""
}

func NewMetric(name, mtype, value string) (*Metric, error) {
	m := new(Metric)
	m.ID = name
	m.MType = mtype

	switch m.MType {
	case domain.MetricGauge:
		if value == "" {
			return m, nil
		}
		cv, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return m, domain.ErrInvalidData
		}
		m.Value = &cv
	case domain.MetricCounter:
		if value == "" {
			return m, nil
		}
		iv, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return m, domain.ErrInvalidData
		}
		m.Delta = &iv
	}

	return m, nil
}

func (m *Metric) Validate() error {
	if m.ID == "" {
		return domain.ErrWrongMetricName
	}

	if m.MType == "" {
		return domain.ErrWrongMetricType
	}

	switch m.MType {
	case domain.MetricGauge, domain.MetricCounter:
		break
	default:
		return domain.ErrWrongMetricType
	}

	if m.GetValue() == "" {
		return domain.ErrEmptyMetricValue
	}

	return nil
}
