package metric

import (
	"encoding/binary"
	"strconv"

	"github.com/konstantin-kukharev/metrics/internal"
)

// Counter Тип int64
type Counter struct{}

func (c *Counter) GetName() string {
	return internal.MetricCounter
}

func (c *Counter) Encode(v string) ([]byte, error) {
	iv, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return []byte{}, internal.ErrInvalidData
	}

	bv := make([]byte, 8)
	binary.LittleEndian.PutUint64(bv, uint64(iv))
	return bv, nil
}

func (c *Counter) Decode(v []byte) (string, error) {
	iv := int64(binary.LittleEndian.Uint64(v))
	sv := strconv.FormatInt(iv, 10)
	return sv, nil
}

// Aggregate новое значение должно добавляться к предыдущему,
// если какое-то значение уже было известно серверу.
func (c *Counter) Aggregate(v ...[]byte) ([]byte, error) {
	var summ int64
	for _, b := range v {
		iv := int64(binary.LittleEndian.Uint64(b))
		summ += iv
	}

	bv := make([]byte, 8)
	binary.LittleEndian.PutUint64(bv, uint64(summ))

	return bv, nil
}
