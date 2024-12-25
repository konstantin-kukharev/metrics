package transport

import (
	"context"

	"github.com/konstantin-kukharev/metrics/domain/event"
)

type EventBus struct {
	addMetricBus chan event.MetricAdd
	addListeners []chan<- event.MetricAdd
}

func (eb *EventBus) AddMetric(e event.MetricAdd) {
	if eb.addMetricBus == nil || len(eb.addListeners) == 0 {
		return
	}

	eb.addMetricBus <- e
}

func (eb *EventBus) AddListener(l ...chan<- event.MetricAdd) {
	eb.addListeners = append(eb.addListeners, l...)
}

func (eb *EventBus) Run(ctx context.Context) error {
	for {
		select {
		case e := <-eb.addMetricBus:
			for _, l := range eb.addListeners {
				l <- e
			}
		case <-ctx.Done():
			for _, l := range eb.addListeners {
				close(l)
			}
			eb.addListeners = []chan<- event.MetricAdd{}
			close(eb.addMetricBus)
			eb.addMetricBus = nil

			return nil
		}
	}
}

func NewEventBus(c chan event.MetricAdd) *EventBus {
	eb := new(EventBus)
	eb.addMetricBus = c

	return eb
}
