package speed

import "github.com/miaokobot/miaospeed/interfaces"

type Speed struct {
	AvgSpeed  uint64
	MaxSpeed  uint64
	TotalSize uint64
	Speeds    []uint64
}

func (m *Speed) Type() interfaces.SlaveRequestMacroType {
	return interfaces.MacroSpeed
}

func (m *Speed) Run(proxy interfaces.Vendor, r *interfaces.SlaveRequest) error {
	Once(m, proxy, &r.Configs)

	return nil
}
