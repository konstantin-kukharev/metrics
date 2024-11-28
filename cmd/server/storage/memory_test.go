package storage

import "testing"

func TestMemoryStorage(t *testing.T) {
	type fields struct {
		key   string
		value []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "set - get test",
			fields: fields{
				key:   "t1",
				value: []byte("test 1"),
			},
			want: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := NewMemStorage()
			ms.Set(
				tt.fields.key,
				tt.fields.value,
			)

			val, ok := ms.Get(
				tt.fields.key,
			)

			if !ok {
				t.Error("err1")
			}

			if got := string(val); got != tt.want {
				t.Errorf("value = %v, want %v", got, tt.want)
			}
		})
	}
}
