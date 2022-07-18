package persecondspeed

import (
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/service/macros/speed"
)

type PerSecondSpeed struct {
	interfaces.PerSecondSpeedDS
}

func (m *PerSecondSpeed) Type() interfaces.SlaveRequestMatrixType {
	return interfaces.MatrixPerSecondSpeed
}

func (m *PerSecondSpeed) MacroJob() interfaces.SlaveRequestMacroType {
	return interfaces.MacroSpeed
}

func (m *PerSecondSpeed) Extract(entry interfaces.SlaveRequestMatrixEntry, macro interfaces.SlaveRequestMacro) {
	if mac, ok := macro.(*speed.Speed); ok {
		m.Speeds = mac.Speeds[:]
		m.Average = mac.AvgSpeed
		m.Max = mac.MaxSpeed
	}
}
