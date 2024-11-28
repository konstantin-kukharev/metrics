package metric

import (
	"strconv"

	"github.com/konstantin-kukharev/metrics/internal"
)

// Gauge Тип float64 - новое значение должно замещать предыдущее.
func Gauge() *class {
	return &class{
		name: "gauge",
		setter: func(s internal.Storage, k, v string) error {
			if _, err := strconv.ParseFloat(v, 64); err != nil {
				return internal.ErrInvalidData
			}
			s.Set(k, []byte(v))

			return nil
		},
	}
}
