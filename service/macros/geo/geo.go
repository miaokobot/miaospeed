package geo

import (
	"net"
	"strings"

	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils"
)

func RemoteLookup(p interfaces.Vendor, script string, retry int) *interfaces.IPStacks {
	ret := &interfaces.IPStacks{}
	for i := 0; i < retry && ret.Count() == 0; i++ {
		ret = ExecIpCheck(p, script, interfaces.ROptionsTCP)
	}
	return ret
}

type DetectSourceMode int

const (
	DSMDefault DetectSourceMode = iota
	DSMInOnly
	DSMOutOnly
)

func DetectingSource(p interfaces.Vendor, script string, retry int, queryServers []string, mode DetectSourceMode) (in *interfaces.MultiStacks, out *interfaces.MultiStacks) {
	if mode == DSMOutOnly || mode == DSMDefault {
		out = &interfaces.MultiStacks{
			Domain:    p.ProxyInfo().Name,
			IPv4Stack: make([]*interfaces.GeoInfo, 0),
			IPv6Stack: make([]*interfaces.GeoInfo, 0),
			MainStack: &interfaces.GeoInfo{},
		}

		outIpstacks := RemoteLookup(p, script, retry)
		for _, outv6Ip := range outIpstacks.IPv6 {
			if outv6 := RunGeoCheck(p, script, outv6Ip, retry, interfaces.ROptionsTCP6); outv6 != nil {
				out.IPv6Stack = append(out.IPv6Stack, outv6)
				out.MainStack = outv6
			}
		}

		for _, outv4Ip := range outIpstacks.IPv4 {
			if outv4 := RunGeoCheck(p, script, outv4Ip, retry, interfaces.ROptionsTCP); outv4 != nil {
				out.IPv4Stack = append(out.IPv4Stack, outv4)
				out.MainStack = outv4
			}
		}
	}

	if mode == DSMInOnly || mode == DSMDefault {
		inIP := p.ProxyInfo().Address
		if strings.Count(inIP, ":") > 1 {
			// ipv6
			inIP = p.ProxyInfo().Address
		} else {
			// domain
			inIP = strings.Split(p.ProxyInfo().Address, ":")[0]
		}

		domain := inIP
		ipv4 := []string{}
		ipv6 := []string{}
		if net.ParseIP(inIP) == nil {
			ipstacks := utils.LookupIPv46(inIP, retry, queryServers)
			ipv4 = ipstacks.IPv4
			ipv6 = ipstacks.IPv6
		} else {
			if strings.Contains(inIP, ":") {
				ipv6 = []string{inIP}
			} else {
				ipv4 = []string{inIP}
			}
		}

		in = &interfaces.MultiStacks{
			Domain:    domain,
			IPv4Stack: make([]*interfaces.GeoInfo, 0),
			IPv6Stack: make([]*interfaces.GeoInfo, 0),
			MainStack: &interfaces.GeoInfo{},
		}

		for _, ip := range ipv4 {
			in.IPv4Stack = append(in.IPv4Stack, RunGeoCheck(nil, script, ip, retry, "tcp"))
		}
		for _, ip := range ipv6 {
			in.IPv6Stack = append(in.IPv6Stack, RunGeoCheck(nil, script, ip, retry, "tcp"))
		}

		if in.Count() > 0 {
			if len(in.IPv4Stack) > 0 {
				in.MainStack = in.IPv4Stack[0]
			} else {
				in.MainStack = in.IPv6Stack[0]
			}
		}
	}

	return
}
