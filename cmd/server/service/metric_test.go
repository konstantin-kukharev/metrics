package service

import (
	"testing"

	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/internal/metric"
	"github.com/stretchr/testify/assert"
)

func TestMetricServiceSet(t *testing.T) {
	type fields struct {
		t     string
		key   string
		value string
	}

	cfg := settings.NewConfig().WithDebug()
	store := storage.NewMemStorage()
	srv := NewMetric(cfg.Log(), store, &metric.Gauge{}, &metric.Counter{})

	tests := []struct {
		name   string
		fields fields
		srv    Metric
		want   interface{}
	}{
		{
			name: "set test 1",
			fields: fields{
				t:     domain.MetricCounter,
				key:   "test",
				value: "33",
			},
			srv:  srv,
			want: nil,
		},
		{
			name: "set test 2",
			fields: fields{
				t:     domain.MetricCounter,
				key:   "",
				value: "33",
			},
			srv:  srv,
			want: domain.ErrWrongMetricName,
		},
		{
			name: "set test 3",
			fields: fields{
				t:     "hehey",
				key:   "asd",
				value: "33",
			},
			srv:  srv,
			want: domain.ErrWrongMetricType,
		},
		{
			name: "set test 3",
			fields: fields{
				t:     domain.MetricCounter,
				key:   "asd",
				value: "asdasdas",
			},
			srv:  srv,
			want: domain.ErrInvalidData,
		},
		{
			name: "set test 4",
			fields: fields{
				t:     domain.MetricCounter,
				key:   "asd",
				value: "11",
			},
			srv:  srv,
			want: nil,
		},
		{
			name: "set test 5",
			fields: fields{
				t:     domain.MetricCounter,
				key:   "asd",
				value: "11",
			},
			srv:  srv,
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.srv.Set(
				tt.fields.t,
				tt.fields.key,
				tt.fields.value,
			)

			assert.Equal(t, err, tt.want)
		},
		)
	}
}

func TestMetricServiceGet(t *testing.T) {
	cfg := settings.NewConfig().WithDebug()
	ns := storage.NewMemStorage()
	srv := NewMetric(cfg.Log(), ns, &metric.Gauge{}, &metric.Counter{})

	err := srv.Set(
		domain.MetricCounter,
		"addTest",
		"33",
	)

	assert.Nil(t, err)

	assert.Nil(t, err)

	type fields struct {
		t   string
		key string
	}

	type want struct {
		res []byte
		ok  bool
	}

	tests := []struct {
		name   string
		fields fields
		srv    Metric
		want   want
	}{
		{
			name: "get test 1",
			fields: fields{
				t:   domain.MetricCounter,
				key: "test999999",
			},
			srv:  srv,
			want: want{[]byte{}, false},
		},
		{
			name: "get test 2",
			fields: fields{
				t:   domain.MetricCounter,
				key: "addTest",
			},
			srv:  srv,
			want: want{[]byte("33"), true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := tt.srv.Get(
				tt.fields.t,
				tt.fields.key,
			)

			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.res, res)
		},
		)
	}
}
