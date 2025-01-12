package entity

import (
	"database/sql"
	"testing"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	cv := 1.23
	var iv int64 = 123
	tests := []struct {
		name     string
		mtype    string
		value    string
		expect   *Metric
		summ     string
		err      error
		validErr error
	}{
		{
			name:   "valid gauge",
			mtype:  domain.MetricGauge,
			value:  "1.23",
			summ:   "1.23",
			expect: &Metric{ID: "valid gauge", MType: domain.MetricGauge, MValue: MValue{Value: &cv}},
		},
		{
			name:   "valid counter",
			mtype:  domain.MetricCounter,
			value:  "123",
			summ:   "246",
			expect: &Metric{ID: "valid counter", MType: domain.MetricCounter, MValue: MValue{Delta: &iv}},
		},
		{
			name:   "invalid gauge",
			mtype:  domain.MetricGauge,
			value:  "invalid",
			expect: &Metric{ID: "invalid gauge", MType: domain.MetricGauge},
			err:    domain.ErrInvalidData,
		},
		{
			name:   "invalid counter",
			mtype:  domain.MetricCounter,
			value:  "invalid",
			expect: &Metric{ID: "invalid counter", MType: domain.MetricCounter},
			err:    domain.ErrInvalidData,
		},
		{
			name:     "empty gauge",
			mtype:    domain.MetricGauge,
			value:    "",
			summ:     "",
			expect:   &Metric{ID: "empty gauge", MType: domain.MetricGauge, MValue: MValue{Delta: nil, Value: nil}},
			validErr: domain.ErrEmptyMetricValue,
		},
		{
			name:     "empty counter",
			mtype:    domain.MetricCounter,
			value:    "",
			summ:     "",
			expect:   &Metric{ID: "empty counter", MType: domain.MetricCounter, MValue: MValue{Delta: nil, Value: nil}},
			validErr: domain.ErrEmptyMetricValue,
		},
		{
			name:     "wrong Type",
			mtype:    "Wrong",
			value:    "",
			summ:     "",
			expect:   &Metric{ID: "wrong Type", MType: "Wrong", MValue: MValue{Delta: nil, Value: nil}},
			validErr: domain.ErrWrongMetricType,
		},
		{
			name:     "wrong Type",
			mtype:    "",
			value:    "",
			summ:     "",
			expect:   &Metric{ID: "wrong Type", MType: "", MValue: MValue{Delta: nil, Value: nil}},
			validErr: domain.ErrWrongMetricType,
		},
		{
			name:     "",
			mtype:    domain.MetricCounter,
			value:    "",
			summ:     "",
			expect:   &Metric{ID: "", MType: domain.MetricCounter, MValue: MValue{Delta: nil, Value: nil}},
			validErr: domain.ErrWrongMetricName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMetric(tt.name, tt.mtype, tt.value)
			assert.Equal(t, tt.expect, m)
			assert.Equal(t, tt.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.value, m.GetValue())
			m.Aggregate(tt.expect)
			assert.Equal(t, tt.summ, m.GetValue())
			err = m.Validate()
			assert.Equal(t, tt.validErr, err)
		})
	}
}

func TestAddValue(t *testing.T) {
	type increment struct {
		d sql.NullInt64
		f sql.NullFloat64
	}

	var d1 int64 = 123
	v1 := 1.23
	var _ int64 = 123
	tests := []struct {
		name   string
		value  *Metric
		expect *Metric
		inc    increment
		err    error
	}{
		{
			name:   "add gauge",
			value:  &Metric{ID: "gauge1", MType: domain.MetricGauge, MValue: MValue{Value: nil, Delta: nil}},
			expect: &Metric{ID: "gauge1", MType: domain.MetricGauge, MValue: MValue{Value: &v1, Delta: nil}},
			inc:    increment{d: sql.NullInt64{Valid: false, Int64: 0}, f: sql.NullFloat64{Valid: true, Float64: 1.23}},
		},
		{
			name:   "add gauge nil",
			value:  &Metric{ID: "gauge1", MType: domain.MetricGauge, MValue: MValue{Value: nil, Delta: nil}},
			expect: &Metric{ID: "gauge1", MType: domain.MetricGauge, MValue: MValue{Value: nil, Delta: nil}},
			inc:    increment{d: sql.NullInt64{Valid: false, Int64: 0}, f: sql.NullFloat64{Valid: false, Float64: 0}},
		},
		{
			name:   "add counter",
			value:  &Metric{ID: "delta1", MType: domain.MetricCounter, MValue: MValue{Value: nil, Delta: nil}},
			expect: &Metric{ID: "delta1", MType: domain.MetricCounter, MValue: MValue{Value: nil, Delta: &d1}},
			inc:    increment{d: sql.NullInt64{Valid: true, Int64: 123}, f: sql.NullFloat64{Valid: false, Float64: 0}},
		},
		{
			name:   "add counter nil",
			value:  &Metric{ID: "delta1", MType: domain.MetricCounter, MValue: MValue{Value: nil, Delta: nil}},
			expect: &Metric{ID: "delta1", MType: domain.MetricCounter, MValue: MValue{Value: nil, Delta: nil}},
			inc:    increment{d: sql.NullInt64{Valid: false, Int64: 0}, f: sql.NullFloat64{Valid: false, Float64: 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.value.SetValue(tt.inc.d, tt.inc.f)
			assert.Equal(t, tt.expect, tt.value)
		})
	}
}
