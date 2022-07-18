package macros

import (
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils/structs"

	"github.com/miaokobot/miaospeed/service/macros/geo"
	"github.com/miaokobot/miaospeed/service/macros/ping"
	"github.com/miaokobot/miaospeed/service/macros/script"
	"github.com/miaokobot/miaospeed/service/macros/speed"
	"github.com/miaokobot/miaospeed/service/macros/udp"

	"github.com/miaokobot/miaospeed/service/macros/invalid"
)

var registeredList = map[interfaces.SlaveRequestMacroType]func() interfaces.SlaveRequestMacro{
	interfaces.MacroSpeed: func() interfaces.SlaveRequestMacro {
		return &speed.Speed{}
	},
	interfaces.MacroPing: func() interfaces.SlaveRequestMacro {
		return &ping.Ping{}
	},
	interfaces.MacroUDP: func() interfaces.SlaveRequestMacro {
		return &udp.Udp{}
	},
	interfaces.MacroGeo: func() interfaces.SlaveRequestMacro {
		return &geo.Geo{}
	},
	interfaces.MacroScript: func() interfaces.SlaveRequestMacro {
		return &script.Script{}
	},
}

func Find(macroType interfaces.SlaveRequestMacroType) interfaces.SlaveRequestMacro {
	if newFn, ok := registeredList[macroType]; ok && newFn != nil {
		return newFn()
	}

	return &invalid.Invalid{}
}

func FindBatch(macroTypes []interfaces.SlaveRequestMacroType) []interfaces.SlaveRequestMacro {
	return structs.Map(macroTypes, func(m interfaces.SlaveRequestMacroType) interfaces.SlaveRequestMacro {
		return Find(m)
	})
}
