package udp

import (
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
	mapType, filterType := detectNATType(proxy, preconfigs.PROXY_DEFAULT_STUN_SERVER)
	m.NATType = natTypeToString(mapType, filterType)

	return nil
}
