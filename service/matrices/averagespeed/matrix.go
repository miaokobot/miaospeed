package averagespeed

import (
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/service/macros/speed"
)

type AverageSpeed struct {
	interfaces.AverageSpeedDS
}

func (m *AverageSpeed) Type() interfaces.SlaveRequestMatrixType {
	return interfaces.MatrixAverageSpeed
}

func (m *AverageSpeed) MacroJob() interfaces.SlaveRequestMacroType {
	return interfaces.MacroSpeed
}

func (m *AverageSpeed) Extract(entry interfaces.SlaveRequestMatrixEntry, macro interfaces.SlaveRequestMacro) {
	if mac, ok := macro.(*speed.Speed); ok {
		m.Value = mac.AvgSpeed
	}
}
