package repository

import (
	"context"

	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal/graceful"
)

type Metric interface {
	Set(context.Context, ...*entity.Metric) ([]*entity.Metric, error)
	Get(context.Context, *entity.Metric) (*entity.Metric, bool)
	List(context.Context) []*entity.Metric
	graceful.Task
}
