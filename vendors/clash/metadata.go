package clash

import (
	"net/url"

	"github.com/Dreamacro/clash/constant"
)

func urlToMetadata(rawURL string, network constant.NetWork) (addr constant.Metadata, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			return
		}
	}

	addr = constant.Metadata{
		NetWork: network,
		Host:    u.Hostname(),
		DstIP:   nil,
		DstPort: port,
	}
	return
}
