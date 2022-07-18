package geo

import (
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils/structs"
)

type Geo struct {
	InStacks  interfaces.MultiStacks
	OutStacks interfaces.MultiStacks
}

func (m *Geo) Type() interfaces.SlaveRequestMacroType {
	return interfaces.MacroGeo
}

func (m *Geo) Run(proxy interfaces.Vendor, r *interfaces.SlaveRequest) error {
	ipScripts := structs.Filter(r.Configs.Scripts, func(v interfaces.Script) bool {
		return v.Type == interfaces.STypeIP
	})

	ipScript := ""
	if len(ipScripts) > 0 {
		ipScript = ipScripts[0].Content
	}

	inStacks, outStacks := DetectingSource(proxy, ipScript, 3, r.Configs.DNSServers, DSMDefault)
	if inStacks != nil {
		m.InStacks = *inStacks
	}
	if outStacks != nil {
		m.OutStacks = *outStacks
	}

	return nil
}
