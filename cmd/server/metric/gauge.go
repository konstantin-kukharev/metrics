package metric

import (
	"strconv"

	"github.com/konstantin-kukharev/metrics/cmd/server/internal"
)

// Gauge Тип float64 - новое значение должно замещать предыдущее.
func Gauge() Metric {
	return Metric{
		Name: "gauge",
		Setter: func(s internal.Storage, k, v string) error {
			if _, err := strconv.ParseFloat(v, 64); err != nil {
				return internal.ErrInvalidData
			}
			s.Set(k, []byte(v))

			return nil
		},
	}
}
