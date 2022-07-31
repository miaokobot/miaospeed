package local

import (
	"fmt"
	"net/url"
	"strconv"
)

func urlToMetadata(rawURL string) (string, uint16, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", 0, fmt.Errorf("cannot parse the url")
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			return "", 0, fmt.Errorf("unknown url scheme")
		}
	}

	portUint8, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return "", 0, fmt.Errorf("cannot parse the port number")
	}

	return u.Hostname(), uint16(portUint8), nil
}
