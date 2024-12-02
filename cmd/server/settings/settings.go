package settings

type settings struct {
	address string
}

func (s *settings) Address() string {
	return s.address
}

func New() *settings {
	s := &settings{}
	fromFlag(s)

	return s
}
