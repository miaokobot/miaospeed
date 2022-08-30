package script

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils/structs"
	"golang.org/x/sync/semaphore"
)

var scriptControl *semaphore.Weighted

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
	for i := range execScripts {
		script := &execScripts[i]
		go func() {
			scriptControl.Acquire(context.Background(), 1)
			defer scriptControl.Release(1)

			store.Set(script.ID, ExecScript(proxy, script))
			wg.Done()
		}()
	}
	wg.Wait()

	m.Store = store.ForEach()
	return nil
}

func init() {
	// default strict to 32 concurrent script engine
	// can be extended by setting env var
	concurrency, _ := strconv.ParseInt(os.Getenv("MIAOKO_SCRIPT_CONCURRENCY"), 10, 64)
	concurrency = structs.WithInDefault(concurrency, 1, 64, 32)
	scriptControl = semaphore.NewWeighted(concurrency)
}
