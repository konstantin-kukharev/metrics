package dto

type metricTypeValue struct {
	c string
	n string
	v string
}

func (m *metricTypeValue) Type() string {
	return m.c
}

func (m *metricTypeValue) Name() string {
	return m.n
}

func (m *metricTypeValue) Value() string {
	return m.v
}

func NewMetricValue(class, name, value string) *metricTypeValue {
	return &metricTypeValue{c: class, n: name, v: value}
}
