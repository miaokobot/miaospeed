package helpers

import (
	"fmt"
	"reflect"

	"github.com/dop251/goja"
	jsoniter "github.com/json-iterator/go"
)

func VMCheck(v goja.Value) bool {
	return v != nil && !goja.IsNull(v) && !goja.IsUndefined(v)
}

func VMSafeStr(v goja.Value) (string, bool) {
	if VMCheck(v) && v.ExportType().Kind() == reflect.String {
		return v.Export().(string), true
	}
	return "", false
}

func VMSafeBool(v goja.Value) (bool, bool) {
	if VMCheck(v) && v.ExportType().Kind() == reflect.Bool {
		return v.Export().(bool), true
	}
	return false, false
}

func VMSafeInt64(v goja.Value) (int64, bool) {
	if VMCheck(v) && v.ExportType().Kind() == reflect.Int64 {
		return v.Export().(int64), true
	}
	return 0, false
}

func VMSafeObj(vm *goja.Runtime, v goja.Value) (*goja.Object, bool) {
	if VMCheck(v) && v.ExportType().Kind() == reflect.Map {
		vo := v.ToObject(vm)
		if vo != nil {
			return vo, true
		}
	}
	return nil, false
}

func VMSafeMarshal(target interface{}, obj goja.Value, vm *goja.Runtime) error {
	if fn, ok := goja.AssertFunction(vm.Get("__json_stringify")); ok {
		ret, _ := fn(goja.Undefined(), obj)
		if v, ok := VMSafeStr(ret); ok {
			if v == "" {
				return fmt.Errorf("cannot marshal an empty string")
			}
			return jsoniter.UnmarshalFromString(v, target)
		}
		return fmt.Errorf("cannot read from stringify function")
	}
	return fmt.Errorf("cannot find stringify function")
}
