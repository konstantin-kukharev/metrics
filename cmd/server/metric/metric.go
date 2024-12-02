package metric

type class struct {
	name     string
	encoder  func(v string) ([]byte, error)
	decoder  func(v []byte) (string, error)
	addition func(v ...[]byte) ([]byte, error)
}

func (m *class) Name() string {
	return m.name
}

func (m *class) Encode(v string) ([]byte, error) {
	return m.encoder(v)
}

func (m *class) Decode(v []byte) (string, error) {
	return m.decoder(v)
}

func (m *class) Addition(v ...[]byte) ([]byte, error) {
	return m.addition(v...)
}
