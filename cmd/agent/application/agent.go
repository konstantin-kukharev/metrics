package application

import (
	"context"
	"runtime"
	"time"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"go.uber.org/zap"
)

type config interface {
	GetServerAddress() string
	GetPoolInterval() time.Duration
}

type repo interface {
	Set(context.Context, ...*entity.Metric) ([]*entity.Metric, error)
}

type Agent struct {
	log          *logger.Logger
	poolInterval time.Duration
	counter      int64
	updater      repo
}

func (a *Agent) Run(ctx context.Context) error {
	a.log.InfoCtx(ctx, "agent is running")
	for {
		select {
		case <-time.After(a.poolInterval):
			a.log.InfoCtx(ctx, "update pool")
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			c := context.WithoutCancel(ctx)
			_, err := a.updater.Set(c, a.update(&mem)...)
			if err != nil {
				a.log.InfoCtx(ctx, "error while update metrics",
					zap.String("message", err.Error()),
				)
			}
		case <-ctx.Done():
			a.log.InfoCtx(ctx, "agent stopped")

			return nil
		}
	}
}

func (a *Agent) update(mem *runtime.MemStats) []*entity.Metric {
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
		"RandomValue":   internal.RandFloat64(),
	} {
		metric := &entity.Metric{
			ID:    name,
			MType: domain.MetricGauge,
		}
		metric.Value = &val
		list = append(list, metric)
	}

	a.counter += 1
	cnt := &entity.Metric{
		ID:    "PollCount",
		MType: domain.MetricCounter,
	}
	cnt.Delta = &a.counter

	list = append(list, cnt)

	return list
}

func NewAgent(updater repo, app config, l *logger.Logger) *Agent {
	agent := new(Agent)
	agent.poolInterval = app.GetPoolInterval()
	agent.counter = 0
	agent.updater = updater
	agent.log = l

	return agent
}
