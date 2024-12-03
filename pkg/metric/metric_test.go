package metric

import (
	"testing"

	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/stretchr/testify/assert"
)

func TestGaugeEncode(t *testing.T) {
	type fields struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "encode error test",
			fields: fields{
				str: "t1",
			},
			want: internal.ErrInvalidData,
		},
		{
			name: "encode value test",
			fields: fields{
				str: "11",
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x26, 0x40},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Gauge{}
			b, err := g.Encode(tt.fields.str)
			if err != nil {
				assert.Equal(t, tt.want, err)
				return
			}
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestGaugeDecode(t *testing.T) {
	type fields struct {
		b []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "encode value test",
			fields: fields{
				b: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x26, 0x40},
			},
			want: "11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Gauge{}
			b, err := g.Decode(tt.fields.b)
			if err != nil {
				assert.Equal(t, tt.want, err)
				return
			}
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestGaugeAddition(t *testing.T) {
	type fields struct {
		b [][]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "addition value test",
			fields: fields{
				b: [][]byte{
					{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x26, 0x40},
					{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x26, 0x48},
					{0x0, 0x0, 0x0, 0x0, 0x0, 0x48, 0x26, 0x48},
				},
			},
			want: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x26, 0x40},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Gauge{}
			b, err := g.Aggregate(tt.fields.b...)
			if err != nil {
				assert.Equal(t, tt.want, err)
				return
			}
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestCounterAddition(t *testing.T) {
	type fields struct {
		b [][]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "addition value test",
			fields: fields{
				b: [][]byte{
					{0xb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
					{0xb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
					{0xb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
				},
			},
			want: []byte{0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Counter{}
			b, err := g.Aggregate(tt.fields.b...)
			if err != nil {
				assert.Equal(t, tt.want, err)
				return
			}
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestCounterEncode(t *testing.T) {
	type fields struct {
		str string
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "encode error test",
			fields: fields{
				str: "t1",
			},
			want: internal.ErrInvalidData,
		},
		{
			name: "encode value test",
			fields: fields{
				str: "33",
			},
			want: []byte{0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Counter{}
			b, err := g.Encode(tt.fields.str)
			if err != nil {
				assert.Equal(t, tt.want, err)
				return
			}
			assert.Equal(t, tt.want, b)
		})
	}
}

func TestCounterDecode(t *testing.T) {
	type fields struct {
		b []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "decode value test",
			fields: fields{
				b: []byte{0xb, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			want: "11",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := Counter{}
			b, err := g.Decode(tt.fields.b)
			if err != nil {
				assert.Equal(t, tt.want, err)
				return
			}
			assert.Equal(t, tt.want, b)
			assert.Equal(t, internal.MetricCounter, g.GetName())
		})
	}
}
