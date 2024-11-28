package metric

import "github.com/konstantin-kukharev/metrics/internal"

type class struct {
	name   string
	setter func(s internal.Storage, k, v string) error
}

func (m *class) Setter() func(s internal.Storage, k, v string) error {
	return m.setter
}

func (m *class) Name() string {
	return m.name
}
