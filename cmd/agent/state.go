package main

import (
	"math/rand/v2"
	"runtime"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type RuntimeMetric struct {
	counter int64
}

func (mr *RuntimeMetric) List(mem *runtime.MemStats) []*entity.Metric {
	list := make([]*entity.Metric, 0)

	for name, val := range map[string]float64{
		"Alloc":         float64(mem.Alloc),
		"BuckHashSys":   float64(mem.BuckHashSys),
		"Frees":         float64(mem.Frees),
		"GCCPUFraction": float64(mem.GCCPUFraction),
		"GCSys":         float64(mem.GCSys),
		"HeapAlloc":     float64(mem.HeapAlloc),
		"HeapIdle":      float64(mem.HeapIdle),
		"HeapInuse":     float64(mem.HeapInuse),
		"HeapObjects":   float64(mem.HeapObjects),
		"HeapReleased":  float64(mem.HeapReleased),
		"HeapSys":       float64(mem.HeapSys),
		"LastGC":        float64(mem.LastGC),
		"Lookups":       float64(mem.Lookups),
		"MCacheInuse":   float64(mem.MCacheInuse),
		"MCacheSys":     float64(mem.MCacheSys),
		"MSpanInuse":    float64(mem.MSpanInuse),
		"MSpanSys":      float64(mem.MSpanSys),
		"Mallocs":       float64(mem.Mallocs),
		"NextGC":        float64(mem.NextGC),
		"NumForcedGC":   float64(mem.NumForcedGC),
		"NumGC":         float64(mem.NumGC),
		"OtherSys":      float64(mem.OtherSys),
		"PauseTotalNs":  float64(mem.PauseTotalNs),
		"StackInuse":    float64(mem.StackInuse),
		"StackSys":      float64(mem.StackSys),
		"Sys":           float64(mem.Sys),
		"TotalAlloc":    float64(mem.TotalAlloc),
		"RandomValue":   rand.Float64(),
	} {
		metric := &entity.Metric{
			ID:    name,
			MType: domain.MetricGauge,
		}
		metric.Value = &val
		list = append(list, metric)
	}

	mr.counter += 1
	cnt := &entity.Metric{
		ID:    "PollCount",
		MType: domain.MetricCounter,
	}
	cnt.Delta = &mr.counter

	list = append(list, cnt)

	return list
}

func NewRuntimeMetric() *RuntimeMetric {
	ms := new(RuntimeMetric)
	ms.counter = 0

	return ms
}
