package geo

import (
	"net"
	"strings"

	"github.com/miaokobot/miaospeed/engine"
	"github.com/miaokobot/miaospeed/engine/helpers"
	"github.com/miaokobot/miaospeed/interfaces"
)

func ExecIpCheck(p interfaces.Vendor, script string, network interfaces.RequestOptionsNetwork) (ipstacks *interfaces.IPStacks) {
	ipstacks = (&interfaces.IPStacks{}).Init()

	vm := engine.VMNewWithVendor(p, network)
	vm.RunString(engine.PREDEFINED_SCRIPT + engine.DEFAULT_IP_SCRIPT + script)
	caller := "ip_resolve_default"
	if engine.HasFunction(vm, "ip_resolve") {
		caller = "ip_resolve"
	}

	ret, err := engine.ExecTaskCallback(vm, caller)
	if engine.ThrowExecTaskErr("IPResolve", err) {
		return
	} else {
		ipQuery := []string{}
		helpers.VMSafeMarshal(&ipQuery, ret, vm)
		for _, ip := range ipQuery {
			if net.ParseIP(ip) != nil {
				if !strings.Contains(ip, ":") {
					ipstacks.IPv4 = append(ipstacks.IPv4, ip)
				} else {
					ipstacks.IPv6 = append(ipstacks.IPv6, ip)
				}
			}
		}
	}

	return
}

func ExecGeoCheck(p interfaces.Vendor, script string, ip string, network interfaces.RequestOptionsNetwork) *interfaces.GeoInfo {
	vm := engine.VMNewWithVendor(p, network)
	if script == "" {
		script = engine.DEFAULT_GEOIP_SCRIPT
	}
	vm.RunString(engine.PREDEFINED_SCRIPT + script)

	ret, err := engine.ExecTaskCallback(vm, "handler", ip)
	if engine.ThrowExecTaskErr("GeoCheck", err) {
		return nil
	} else {
		geoInfo := &interfaces.GeoInfo{}
		if err := helpers.VMSafeMarshal(geoInfo, ret, vm); err == nil {
			return geoInfo
		}
	}

	return nil
}
