package metric

import (
	"encoding/binary"
	"math"
	"strconv"

	"github.com/konstantin-kukharev/metrics/internal"
)

// Gauge Тип float64 floatPrecision = internal.DefaultFloatPrecision
type Gauge struct{}

func (g *Gauge) GetName() string {
	return internal.MetricGauge
}

func (g *Gauge) Encode(v string) ([]byte, error) {
	cv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return []byte{}, internal.ErrInvalidData
	}

	bits := math.Float64bits(cv)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes, nil
}

func (g *Gauge) Decode(v []byte) (string, error) {
	bits := binary.LittleEndian.Uint64(v)
	fv := math.Float64frombits(bits)
	sv := strconv.FormatFloat(fv, 'f', internal.DefaultFloatPrecision, 64)
	return sv, nil
}

// Aggregate новое значение должно замещать предыдущее.
func (g *Gauge) Aggregate(v ...[]byte) ([]byte, error) {
	return v[0], nil
}
