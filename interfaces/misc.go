package interfaces

type RequestOptionsNetwork string

const (
	ROptionsTCP  RequestOptionsNetwork = "tcp"
	ROptionsTCP6 RequestOptionsNetwork = "tcp6"
)

func (ron *RequestOptionsNetwork) String() string {
	if ron == nil {
		return "tcp"
	}

	switch *ron {
	case ROptionsTCP:
		return "tcp"
	case ROptionsTCP6:
		return "tcp6"
	}

	return "tcp"
}

type RequestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Cookies map[string]string
	Body    []byte
	NoRedir bool
	Network RequestOptionsNetwork
}
