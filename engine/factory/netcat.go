package factory

import (
	"github.com/dop251/goja"
	"github.com/miaokobot/miaospeed/engine/helpers"
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/vendors"
)

func NetCatFactory(vm *goja.Runtime, p interfaces.Vendor, network interfaces.RequestOptionsNetwork) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		addr, _ := helpers.VMSafeStr(call.Argument(0))
		data, _ := helpers.VMSafeStr(call.Argument(1))
		params, _ := helpers.VMSafeObj(vm, call.Argument(2))

		retry := 0
		useHost := false
		timeout := int64(3000)

		if params != nil {
			if v, ok := helpers.VMSafeBool(params.Get("useHost")); ok {
				useHost = v
			}
			if v, ok := helpers.VMSafeInt64(params.Get("timeout")); ok {
				timeout = v
			}
			if v, ok := helpers.VMSafeInt64(params.Get("retry")); ok {
				retry = int(v)
			}
		}

		if useHost {
			p = nil
		}

		returns, err := vendors.NetCatWithRetry(p, retry, timeout, addr, []byte(data), network)

		retMap := map[string]string{
			"error": "",
			"data":  string(returns),
		}
		if err != nil {
			retMap["error"] = err.Error()
		}

		return vm.ToValue(retMap)
	}
}
