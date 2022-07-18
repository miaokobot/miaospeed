package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/console"
	"github.com/dop251/goja_nodejs/require"

	"github.com/miaokobot/miaospeed/engine/factory"
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils"
	"github.com/miaokobot/miaospeed/utils/structs"
)

func VMNew() *goja.Runtime {
	rt := goja.New()
	new(require.Registry).Enable(rt)
	console.Enable(rt)

	rt.Set("print", factory.PrintFactory(rt, "Script Print |", utils.LTInfo))
	rt.Set("debug", factory.PrintFactory(rt, "Script Debug |", utils.LTLog))

	rt.SetMaxCallStackSize(1024)
	return rt
}

func VMNewWithVendor(p interfaces.Vendor, network interfaces.RequestOptionsNetwork) *goja.Runtime {
	vm := VMNew()

	vm.Set("fetch", factory.FetchFactory(vm, p, network))
	vm.Set("netcat", factory.NetCatFactory(vm, p, network))

	return vm
}

func IsNotExtractError(err error) bool {
	if err != nil {
		return err.Error() == "cannot extract function from vm"
	}
	return false
}

func ThrowExecTaskErr(scenario string, err error) bool {
	if err != nil {
		if !IsNotExtractError(err) {
			utils.DErrorf("Engine Error | scenario=%s error=%s", scenario, err.Error())
		}
		return true
	}
	return false
}

func HasFunction(vm *goja.Runtime, caller string) bool {
	_, ok := goja.AssertFunction(vm.Get(caller))
	return ok
}

func ExecTaskCallback(vm *goja.Runtime, caller string, args ...interface{}) (ret goja.Value, err error) {
	utils.WrapError("Exec task callback error", func() error {
		if vm == nil {
			ret, err = goja.Undefined(), fmt.Errorf("vm is not initialized")
			return nil
		}

		fn, ok := goja.AssertFunction(vm.Get(caller))
		if !ok {
			ret, err = goja.Undefined(), fmt.Errorf("cannot extract function from vm")
			return nil
		}

		values := []goja.Value{}
		for _, arg := range args {
			values = append(values, vm.ToValue(arg))
		}

		ret, err = fn(goja.Undefined(), values...)
		return nil
	})

	return
}

func RunWithTimeout(vm *goja.Runtime, timeout time.Duration, fn func() (goja.Value, error)) (ret goja.Value, err error) {
	vmLock := sync.Mutex{}
	finished := false

	if timeout > 0 {
		timeout = structs.WithIn(timeout, time.Second, time.Minute)
		time.AfterFunc(timeout, func() {
			vmLock.Lock()
			defer vmLock.Unlock()

			if !finished {
				finished = true
				vm.Interrupt("script executing too long")
			}
		})

	}

	ret, err = fn()

	vmLock.Lock()
	finished = true
	vmLock.Unlock()

	return
}
