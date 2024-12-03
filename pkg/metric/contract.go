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
