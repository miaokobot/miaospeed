package factory

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/miaokobot/miaospeed/utils"
)

func PrintFactory(vm *goja.Runtime, prefix string, logType utils.LogType) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		prep := prefix
		args := call.Arguments
		pass := make([]interface{}, len(args))
		for i := 0; i < len(args); i++ {
			prep += " %v"
			pass[i] = args[i]
		}
		prep = fmt.Sprintf(prep, pass...)
		utils.DBase(logType, prep)
		return vm.ToValue(prep)
	}
}
