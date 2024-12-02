package metric

import (
	"encoding/binary"
	"strconv"

	"github.com/konstantin-kukharev/metrics/internal"
)

// Counter Тип int64 - новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func Counter() *class {
	return &class{
		name: internal.MetricCounter,
		encoder: func(v string) ([]byte, error) {
			iv, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return []byte{}, internal.ErrInvalidData
			}

			bv := make([]byte, 8)
			binary.LittleEndian.PutUint64(bv, uint64(iv))
			return bv, nil
		},
		decoder: func(v []byte) (string, error) {
			iv := int64(binary.LittleEndian.Uint64(v))
			sv := strconv.FormatInt(iv, 10)
			return sv, nil
		},
		addition: func(v ...[]byte) ([]byte, error) {
			var summ int64
			for _, b := range v {
				iv := int64(binary.LittleEndian.Uint64(b))
				summ += iv
			}

			bv := make([]byte, 8)
			binary.LittleEndian.PutUint64(bv, uint64(summ))

			return bv, nil
		},
	}
}
