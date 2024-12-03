package state

import (
	"runtime"

	"github.com/konstantin-kukharev/metrics/pkg/metric"
)

type Memory interface {
	Update(m *runtime.MemStats)
	Data() []metric.Value
}
