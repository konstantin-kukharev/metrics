package entity

import (
	"testing"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewMetric(t *testing.T) {
	cv := 1.23
	var iv int64 = 123
	tests := []struct {
		name   string
		mtype  string
		value  string
		expect *Metric
		err    error
	}{
		{
			name:   "valid gauge",
			mtype:  domain.MetricGauge,
			value:  "1.23",
			expect: &Metric{ID: "valid gauge", MType: domain.MetricGauge, Value: &cv},
		},
		{
			name:   "valid counter",
			mtype:  domain.MetricCounter,
			value:  "123",
			expect: &Metric{ID: "valid counter", MType: domain.MetricCounter, Delta: &iv},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMetric(tt.name, tt.mtype, tt.value)
			assert.Equal(t, tt.expect, m)
			assert.Equal(t, tt.err, err)
		})
	}
}
