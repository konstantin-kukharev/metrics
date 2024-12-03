package storage

import (
	"fmt"
	"testing"

	"github.com/konstantin-kukharev/metrics/pkg/metric"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	type fields struct {
		t     metric.Type
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "set - get test",
			fields: fields{
				t:     &metric.Counter{},
				key:   "test",
				value: "33",
			},
			want: "33",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				ms := NewMemStorage()
				err := ms.Set(
					tt.fields.t,
					tt.fields.key,
					tt.fields.value,
				)

				if err != nil {
					t.Error(err.Error())
				}

				val, ok := ms.Get(
					tt.fields.t,
					tt.fields.key,
				)

				fmt.Println(ms.storage)

				if !ok {
					t.Error("err1")
				}

				if val != tt.want {
					t.Errorf("value = %v, want %v", val, tt.want)
				}

				data := ms.List()
				assert.Equal(t, 1, len(data))
				got, err := metric.NewValue(tt.fields.t, tt.fields.key, val)
				assert.Nil(t, err)
				assert.Equal(
					t,
					got,
					data[0],
				)
			},
		)
	}
}
