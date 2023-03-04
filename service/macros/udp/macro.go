package udp

import (
	"strings"

	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/preconfigs"
)

type Udp struct {
	NATType string
}

func (m *Udp) Type() interfaces.SlaveRequestMacroType {
	return interfaces.MacroUDP
}

func (m *Udp) Run(proxy interfaces.Vendor, r *interfaces.SlaveRequest) error {
	stunURL := strings.TrimSpace(r.Configs.STUNURL)
	if stunURL == "" {
		stunURL = preconfigs.PROXY_DEFAULT_STUN_SERVER
	}

	mapType, filterType := detectNATType(proxy, stunURL)
	m.NATType = natTypeToString(mapType, filterType)

	return nil
}
