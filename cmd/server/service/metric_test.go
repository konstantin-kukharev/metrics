package service

import (
	"testing"

	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/konstantin-kukharev/metrics/pkg/metric"
	"github.com/stretchr/testify/assert"
)

func newService() Metric {
	cfg := &settings.Config{Address: internal.DefaultServerAddr}
	store := storage.NewMemStorage()
	return NewMetric(cfg, store, &metric.Gauge{}, &metric.Counter{})
}

func TestMetricServiceSet(t *testing.T) {

	type fields struct {
		t     string
		key   string
		value string
	}

	tests := []struct {
		name   string
		fields fields
		srv    Metric
		want   interface{}
	}{
		{
			name: "set test 1",
			fields: fields{
				t:     internal.MetricCounter,
				key:   "test",
				value: "33",
			},
			srv:  newService(),
			want: nil,
		},
		{
			name: "set test 2",
			fields: fields{
				t:     internal.MetricCounter,
				key:   "",
				value: "33",
			},
			srv:  newService(),
			want: internal.ErrWrongMetricName,
		},
		{
			name: "set test 3",
			fields: fields{
				t:     "hehey",
				key:   "asd",
				value: "33",
			},
			srv:  newService(),
			want: internal.ErrWrongMetricType,
		},
		{
			name: "set test 3",
			fields: fields{
				t:     internal.MetricCounter,
				key:   "asd",
				value: "asdasdas",
			},
			srv:  newService(),
			want: internal.ErrInvalidData,
		},
		{
			name: "set test 4",
			fields: fields{
				t:     internal.MetricCounter,
				key:   "asd",
				value: "11",
			},
			srv:  newService(),
			want: nil,
		},
		{
			name: "set test 5",
			fields: fields{
				t:     internal.MetricCounter,
				key:   "asd",
				value: "11",
			},
			srv:  newService(),
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
	ns := newService()
	err := ns.Set(
		internal.MetricCounter,
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
				t:   internal.MetricCounter,
				key: "test999999",
			},
			srv:  ns,
			want: want{[]byte{}, false},
		},
		{
			name: "get test 2",
			fields: fields{
				t:   internal.MetricCounter,
				key: "addTest",
			},
			srv:  ns,
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
