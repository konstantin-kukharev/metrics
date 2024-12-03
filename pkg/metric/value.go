package metric

type typedValue struct {
	c Type
	n string
	v []byte
}

func (m *typedValue) Type() string {
	return m.c.GetName()
}

func (m *typedValue) Name() string {
	return m.n
}

func (m *typedValue) GetValue() (string, error) {
	return m.c.Decode(m.v)
}

func NewValue(t Type, name, value string) (*typedValue, error) {
	val, err := t.Encode(value)
	if err != nil {
		return nil, err
	}

	return &typedValue{c: t, n: name, v: val}, nil
}
