package metric

import (
	"encoding/binary"
	"strconv"

	"github.com/konstantin-kukharev/metrics/cmd/server/internal"
)

// Counter Тип int64 - новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func Counter() Metric {
	return Metric{
		Name: "counter",
		Setter: func(s internal.Storage, k, v string) error {
			var summ int64
			i1, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return internal.ErrInvalidData
			}

			summ += i1

			val, ok := s.Get(k)
			if ok {
				i2 := int64(binary.LittleEndian.Uint64(val))
				summ += i2
			}

			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(summ))
			s.Set(k, []byte(v))

			return nil
		},
	}
}
