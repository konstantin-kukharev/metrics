package memory

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal/graceful"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"go.uber.org/zap"
)

func TestRepo(t *testing.T) {
	var i int64 = 123
	var f float64 = 1.23

	ctx, cncl := context.WithCancel(context.Background())
	l, err := logger.NewLogger(zap.DebugLevel)
	if err != nil {
		log.Panic(err)
	}
	ctx = l.WithContextFields(ctx,
		zap.String("app", "repo test"))
	defer l.Sync()

	repos := NewMetric(l)

	gs := graceful.NewGracefulShutdown(ctx, 1*time.Second)
	gs.AddTask(repos)

	time.Sleep(1 * time.Second)

	tests := []struct {
		name        string
		update      *entity.Metric
		getExpected []*entity.Metric
		err         error
	}{
		{
			name:   "add gauge",
			update: &entity.Metric{ID: "gauge", MType: domain.MetricGauge, MValue: entity.MValue{Value: &f, Delta: nil}},
			getExpected: []*entity.Metric{
				{ID: "gauge", MType: domain.MetricGauge, MValue: entity.MValue{Value: &f, Delta: nil}},
			},
		},
		{
			name:   "add counter",
			update: &entity.Metric{ID: "counter", MType: domain.MetricCounter, MValue: entity.MValue{Value: nil, Delta: &i}},
			getExpected: []*entity.Metric{
				{ID: "counter", MType: domain.MetricCounter, MValue: entity.MValue{Value: nil, Delta: &i}},
			},
		},
	}

	expextedList := make([]*entity.Metric, 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins, err := repos.Set(ctx, tt.update)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.getExpected, ins)
			expextedList = append(expextedList, ins...)
			get, ok := repos.Get(ctx, ins...)
			assert.Equal(t, true, ok)
			assert.Equal(t, tt.getExpected, get)
			list := repos.List(ctx)
			assert.Equal(t, expextedList, list)
		})
	}

	cncl()
}
