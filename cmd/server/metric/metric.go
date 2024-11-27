package metric

import "github.com/konstantin-kukharev/metrics/cmd/server/internal"

type Metric struct {
	Name   string
	Setter func(s internal.Storage, k, v string) error
}
