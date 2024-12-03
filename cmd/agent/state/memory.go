package state

import (
	"math/rand/v2"
	"runtime"
	"strconv"
	"sync"

	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/konstantin-kukharev/metrics/pkg/metric"
)

type memory struct {
	mx      *sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

func NewMemory() *memory {
	return &memory{
		mx:      &sync.RWMutex{},
		gauge:   map[string]float64{},
		counter: map[string]int64{"PollCount": 0},
	}
}

func (d *memory) Update(m *runtime.MemStats) {
	runtime.ReadMemStats(m)
	d.mx.Lock()
	defer d.mx.Unlock()

	d.gauge = map[string]float64{
		"Alloc":         float64(m.Alloc),
		"BuckHashSys":   float64(m.BuckHashSys),
		"Frees":         float64(m.Frees),
		"GCCPUFraction": float64(m.GCCPUFraction),
		"GCSys":         float64(m.GCSys),
		"HeapAlloc":     float64(m.HeapAlloc),
		"HeapIdle":      float64(m.HeapIdle),
		"HeapInuse":     float64(m.HeapInuse),
		"HeapObjects":   float64(m.HeapObjects),
		"HeapReleased":  float64(m.HeapReleased),
		"HeapSys":       float64(m.HeapSys),
		"LastGC":        float64(m.LastGC),
		"Lookups":       float64(m.Lookups),
		"MCacheInuse":   float64(m.MCacheInuse),
		"MCacheSys":     float64(m.MCacheSys),
		"MSpanInuse":    float64(m.MSpanInuse),
		"MSpanSys":      float64(m.MSpanSys),
		"Mallocs":       float64(m.Mallocs),
		"NextGC":        float64(m.NextGC),
		"NumForcedGC":   float64(m.NumForcedGC),
		"NumGC":         float64(m.NumGC),
		"OtherSys":      float64(m.OtherSys),
		"PauseTotalNs":  float64(m.PauseTotalNs),
		"StackInuse":    float64(m.StackInuse),
		"StackSys":      float64(m.StackSys),
		"Sys":           float64(m.Sys),
		"TotalAlloc":    float64(m.TotalAlloc),
	}

	d.gauge["RandomValue"] = rand.Float64()
	d.counter["PollCount"] = d.counter["PollCount"] + 1
}

func (d *memory) Data() []metric.Value {
	l := len(d.counter) + len(d.gauge)
	res := make([]metric.Value, 0, l)
	d.mx.RLock()
	defer d.mx.RUnlock()

	for n, v := range d.gauge {
		val, err := metric.NewValue(
			&metric.Gauge{},
			n,
			strconv.FormatFloat(v, 'f', internal.DefaultFloatPrecision, 64),
		)
		//todo: catch err
		if err != nil {
			continue
		}

		res = append(
			res,
			val,
		)
	}

	for n, v := range d.counter {
		val, err := metric.NewValue(
			&metric.Counter{},
			n,
			strconv.FormatInt(v, 10),
		)
		//todo: catch err
		if err != nil {
			continue
		}

		res = append(
			res,
			val,
		)
	}

	return res
}
