package script

import (
	"sync"

	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils/structs"
)

type Script struct {
	Store map[string]interfaces.ScriptResult
}

func (m *Script) Type() interfaces.SlaveRequestMacroType {
	return interfaces.MacroScript
}

func (m *Script) Run(proxy interfaces.Vendor, r *interfaces.SlaveRequest) error {
	store := structs.NewAsyncMap[string, interfaces.ScriptResult]()
	execScripts := structs.Filter(r.Configs.Scripts, func(v interfaces.Script) bool {
		return v.Type == interfaces.STypeMedia
	})

	wg := sync.WaitGroup{}
	wg.Add(len(execScripts))
	for _, script := range execScripts {
		store.Set(script.ID, ExecScript(proxy, &script))
	}

	m.Store = store.ForEach()
	return nil
}
