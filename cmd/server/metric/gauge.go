package metric

import (
	"encoding/binary"
	"math"
	"strconv"

	"github.com/konstantin-kukharev/metrics/internal"
)

// Gauge Тип float64 - новое значение должно замещать предыдущее.
func Gauge() *class {
	return &class{
		name: internal.MetricGauge,
		encoder: func(v string) ([]byte, error) {
			cv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return []byte{}, internal.ErrInvalidData
			}

			bits := math.Float64bits(cv)
			bytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(bytes, bits)
			return bytes, nil
		},
		decoder: func(v []byte) (string, error) {
			bits := binary.LittleEndian.Uint64(v)
			fv := math.Float64frombits(bits)
			sv := strconv.FormatFloat(fv, 'f', 3, 64)
			return sv, nil
		},
		addition: func(v ...[]byte) ([]byte, error) {
			return v[0], nil
		},
	}
}
