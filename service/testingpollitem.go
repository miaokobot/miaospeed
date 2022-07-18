package service

import (
	"sync"
	"time"

	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/service/macros"
	"github.com/miaokobot/miaospeed/service/macros/invalid"
	"github.com/miaokobot/miaospeed/service/matrices"
	"github.com/miaokobot/miaospeed/service/taskpoll"
	"github.com/miaokobot/miaospeed/utils"
	"github.com/miaokobot/miaospeed/utils/structs"
	"github.com/miaokobot/miaospeed/vendors"
)

type TestingPollItem struct {
	id   string
	name string

	request  *interfaces.SlaveRequest
	matrices []interfaces.SlaveRequestMatrixEntry
	macros   []interfaces.SlaveRequestMacroType
	results  *structs.AsyncArr[interfaces.SlaveEntrySlot]

	onProcess func(self *TestingPollItem, idx int, result interfaces.SlaveEntrySlot)
	onExit    func(self *TestingPollItem, exitCode taskpoll.TaskPollExitCode)

	onProcessLock sync.Mutex
	exitOnce      sync.Once
}

func (tpi *TestingPollItem) ID() string {
	return tpi.id
}

func (tpi *TestingPollItem) TaskName() string {
	return tpi.name
}

func (tpi *TestingPollItem) Weight() uint {
	// TODO: could arrange weight based on task size
	// or customized rules

	return 1
}

func (tpi *TestingPollItem) Count() int {
	return len(tpi.request.Nodes)
}

func (tpi *TestingPollItem) Yield(idx int, tpc *taskpoll.TaskPollController) {
	result := interfaces.SlaveEntrySlot{
		ProxyInfo:      interfaces.ProxyInfo{},
		InvokeDuration: -1,
		Matrices:       []interfaces.MatrixResponse{},
	}

	defer func() {
		utils.WrapErrorPure("Task yield error", recover())

		tpi.results.Push(result)

		tpi.onProcessLock.Lock()
		defer tpi.onProcessLock.Unlock()
		tpi.onProcess(tpi, idx, result)
	}()

	node := tpi.request.Nodes[idx]
	vendor := vendors.Find(tpi.request.Vendor).Build(node.Name, node.Payload)
	result.ProxyInfo = vendor.ProxyInfo()
	macroMap := structs.NewAsyncMap[interfaces.SlaveRequestMacroType, interfaces.SlaveRequestMacro]()

	startTime := time.Now().UnixMilli()
	wg := sync.WaitGroup{}
	wg.Add(len(tpi.macros))
	for _, macro := range tpi.macros {
		macroName := macro
		go func() {
			macro := macros.Find(macroName)
			macro.Run(vendor, tpi.request)
			macroMap.Set(macroName, macro)
			wg.Done()
		}()
	}
	wg.Wait()
	endTime := time.Now().UnixMilli()
	result.InvokeDuration = endTime - startTime

	result.Matrices = structs.Map(tpi.matrices, func(me interfaces.SlaveRequestMatrixEntry) interfaces.MatrixResponse {
		m := matrices.Find(me.Type)
		macro := macroMap.MustGet(m.MacroJob())
		if macro == nil {
			macro = &invalid.Invalid{}
		}
		m.Extract(me, macro)

		return interfaces.MatrixResponse{
			Type:    m.Type(),
			Payload: utils.ToJSON(m),
		}
	})
}

func (tpi *TestingPollItem) OnExit(exitCode taskpoll.TaskPollExitCode) {
	tpi.exitOnce.Do(func() {
		tpi.onExit(tpi, exitCode)
	})
}

func (tpi *TestingPollItem) Init() taskpoll.TaskPollItem {
	tpi.results = structs.NewAsyncArr[interfaces.SlaveEntrySlot]()
	return tpi
}
