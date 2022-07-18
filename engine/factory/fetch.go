package factory

import (
	"github.com/dop251/goja"
	"github.com/miaokobot/miaospeed/engine/helpers"
	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/vendors"
)

func FetchFactory(vm *goja.Runtime, p interfaces.Vendor, network interfaces.RequestOptionsNetwork) func(call goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		url, _ := helpers.VMSafeStr(call.Argument(0))
		params, _ := helpers.VMSafeObj(vm, call.Argument(1))

		method := "GET"
		body := ""
		useHost := false
		noRedir := false
		retry := 0
		timeout := int64(3000)
		headers := map[string]string{}
		cookies := map[string]string{}

		if params != nil {
			if v, ok := helpers.VMSafeStr(params.Get("method")); ok {
				method = v
			}
			if v, ok := helpers.VMSafeStr(params.Get("body")); ok {
				body = v
			}
			if v, ok := helpers.VMSafeBool(params.Get("useHost")); ok {
				useHost = v
			}
			if v, ok := helpers.VMSafeBool(params.Get("noRedir")); ok {
				noRedir = v
			}
			if v, ok := helpers.VMSafeInt64(params.Get("retry")); ok {
				retry = int(v)
			}
			if v, ok := helpers.VMSafeInt64(params.Get("timeout")); ok {
				timeout = v
			}
			if vo, _ := helpers.VMSafeObj(vm, params.Get("headers")); vo != nil {
				for _, key := range vo.Keys() {
					if vv, ok := helpers.VMSafeStr(vo.Get(key)); ok {
						headers[key] = vv
					}
				}
			}
			if vo, _ := helpers.VMSafeObj(vm, params.Get("cookies")); vo != nil {
				for _, key := range vo.Keys() {
					if vv, ok := helpers.VMSafeStr(vo.Get(key)); ok {
						cookies[key] = vv
					}
				}
			}
		}

		if useHost {
			p = nil
		}
		retBody, resp, redirs := vendors.RequestWithRetry(p, retry, timeout, &interfaces.RequestOptions{
			Method:  method,
			URL:     url,
			Headers: headers,
			Cookies: cookies,
			Body:    []byte(body),
			NoRedir: noRedir,
			Network: network,
		})

		var retMap map[string]interface{} = nil
		if resp != nil {
			retMap = make(map[string]interface{})
			retMap["status"] = resp.Status
			retMap["statusCode"] = resp.StatusCode
			retMap["cookies"] = resp.Cookies()
			retMap["headers"] = resp.Header
			retMap["method"] = method
			retMap["url"] = url
			retMap["body"] = string(retBody)
			retMap["redirects"] = redirs
		}

		if retMap == nil {
			return goja.Null()
		} else {
			return vm.ToValue(retMap)
		}
	}
}
