package script

import (
	"runtime"
	"time"

	"github.com/dop251/goja"
	"github.com/miaokobot/miaospeed/engine"
	"github.com/miaokobot/miaospeed/engine/helpers"
	"github.com/miaokobot/miaospeed/interfaces"
)

func ExecScript(p interfaces.Vendor, script *interfaces.Script) interfaces.ScriptResult {
	s := interfaces.ScriptResult{}
	if script == nil {
		return s
	}

	vm := engine.VMNewWithVendor(p, interfaces.ROptionsTCP)

	startTime := time.Now()
	ret, err := engine.RunWithTimeout(vm, time.Duration(script.TimeoutMillis)*time.Millisecond, func() (goja.Value, error) {
		vm.RunString(engine.PREDEFINED_SCRIPT + script.Content)
		return engine.ExecTaskCallback(vm, "handler")
	})

	s.TimeElapsed = time.Now().UnixMilli() - startTime.UnixMilli()
	if engine.ThrowExecTaskErr("MediaTest", err) {
		// nothing here
	} else if text, ok := helpers.VMSafeStr(ret); ok {
		s.Text = text
	} else if ro, _ := helpers.VMSafeObj(vm, ret); ro != nil {
		if v, ok := helpers.VMSafeStr(ro.Get("text")); ok {
			s.Text = v
		}
		if v, ok := helpers.VMSafeStr(ro.Get("color")); ok {
			s.Color = v
		}
		if v, ok := helpers.VMSafeStr(ro.Get("background")); ok {
			s.Background = v
		}
	}

	vm = nil
	runtime.GC()

	return s
}
