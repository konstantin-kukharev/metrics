package report

import "github.com/konstantin-kukharev/metrics/pkg/metric"

type AgentReporter interface {
	Report(serverAddress string, data []metric.Value) error
}
