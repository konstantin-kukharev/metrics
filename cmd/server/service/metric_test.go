package service

import (
	"testing"

	dto "github.com/konstantin-kukharev/metrics/cmd/server/metric"
	"github.com/konstantin-kukharev/metrics/cmd/server/settings"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/stretchr/testify/assert"
)

func TestMetricServiceSet(t *testing.T) {
	conf := settings.New()
	store := storage.NewMemStorage()
	srv := NewMetric(conf, store, dto.Gauge(), dto.Counter())

	type fields struct {
		t     string
		key   string
		value string
	}

	tests := []struct {
		name   string
		fields fields
		srv    internal.MetricService
		want   interface{}
	}{
		{
			name: "set test 1",
			fields: fields{
				t:     internal.MetricCounter,
				key:   "test",
				value: "33",
			},
			srv:  srv,
			want: nil,
		},
		{
			name: "set test 2",
			fields: fields{
				t:     internal.MetricCounter,
				key:   "",
				value: "33",
			},
			srv:  srv,
			want: internal.ErrWrongMetricName,
		},
		{
			name: "set test 3",
			fields: fields{
				t:     "hehey",
				key:   "asd",
				value: "33",
			},
			srv:  srv,
			want: internal.ErrWrongMetricType,
		},
		{
			name: "set test 3",
			fields: fields{
				t:     internal.MetricCounter,
				key:   "asd",
				value: "asdasdas",
			},
			srv:  srv,
			want: internal.ErrInvalidData,
		},
		{
			name: "set test 4",
			fields: fields{
				t:     internal.MetricCounter,
				key:   "asd",
				value: "11",
			},
			srv:  srv,
			want: nil,
		},
		{
			name: "set test 5",
			fields: fields{
				t:     internal.MetricCounter,
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
	conf := settings.New()
	store := storage.NewMemStorage()
	gauge := dto.Gauge()
	counter := dto.Counter()
	srv := NewMetric(conf, store, gauge, counter)

	err := srv.Set(
		internal.MetricCounter,
		"test123",
		"33",
	)

	assert.Nil(t, err)
	byte3, err := counter.Encode("33")

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
		srv    internal.MetricService
		want   want
	}{
		{
			name: "get test 1",
			fields: fields{
				t:   internal.MetricCounter,
				key: "test1",
			},
			srv:  srv,
			want: want{[]byte{}, false},
		},
		{
			name: "get test 2",
			fields: fields{
				t:   internal.MetricCounter,
				key: "test123",
			},
			srv:  srv,
			want: want{byte3, true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, ok := tt.srv.Get(
				tt.fields.t,
				tt.fields.key,
			)

			assert.Equal(t, ok, tt.want.ok)
			assert.Equal(t, res, tt.want.res)
		},
		)
	}
}
