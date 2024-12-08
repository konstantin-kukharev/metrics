package metric

type Type interface {
	GetName() string
	Encode(value string) ([]byte, error)
	Decode(value []byte) (string, error)
	Aggregate(...[]byte) ([]byte, error)
}

type Value interface {
	Type() string
	Name() string
	GetValue() (string, error)
}

type TypedValue struct {
	c Type
	n string
	v []byte
}

func (m *TypedValue) Type() string {
	return m.c.GetName()
}

func (m *TypedValue) Name() string {
	return m.n
}

func (m *TypedValue) GetValue() (string, error) {
	return m.c.Decode(m.v)
}

func NewValue(t Type, name, value string) (*TypedValue, error) {
	val, err := t.Encode(value)
	if err != nil {
		return nil, err
	}

	return &TypedValue{c: t, n: name, v: val}, nil
}
