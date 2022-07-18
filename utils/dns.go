package utils

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils/structs/memutils"
	"github.com/miaokobot/miaospeed/utils/structs/obliviousmap"
)

var DnsCache *obliviousmap.ObliviousMap[*interfaces.IPStacks]

// queryServer = "8.8.8.8:53"
func DNSLookuper(addr string, queryServers []string) []net.IP {
	if len(queryServers) == 0 {
		result, _ := net.LookupIP(addr)
		return result
	}

	ipSets := map[string]net.IP{}
	for _, server := range queryServers {
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Millisecond * time.Duration(3000),
				}
				return d.DialContext(ctx, network, server)
			},
		}
		addrs, _ := r.LookupIPAddr(context.Background(), addr)
		for _, ia := range addrs {
			ipSets[ia.IP.String()] = ia.IP
		}
	}

	ips := make([]net.IP, len(ipSets))
	j := 0
	for _, ia := range ipSets {
		ips[j] = ia
		j += 1
	}

	return ips
}

func LookupIPv46(addr string, retry int, queryServers []string) *interfaces.IPStacks {
	token := fmt.Sprintf("%v|%v", addr, queryServers)
	if r, ok := DnsCache.Get(token); ok && r != nil {
		return r
	}

	netips := []net.IP{}
	for i := 0; i < retry && len(netips) == 0; i += 1 {
		netips = DNSLookuper(addr, queryServers)
	}
	DLogf("DNS Lookup | dns=%v result=%v", queryServers, netips)

	ipstacks := (&interfaces.IPStacks{}).Init()
	for _, ip := range netips {
		ipStr := ip.String()
		if !strings.Contains(ipStr, ":") {
			ipstacks.IPv4 = append(ipstacks.IPv4, ipStr)
		} else {
			ipstacks.IPv6 = append(ipstacks.IPv6, ipStr)
		}
	}

	if ipstacks.Count() > 0 {
		DnsCache.Set(token, ipstacks)
	} else {
		DWarnf("DNS Resolver | fail to resolve domain=%s", addr)
	}
	return ipstacks
}

func init() {
	memIPStacks := memutils.MemDriverMemory[*interfaces.IPStacks]{}
	memIPStacks.Init()
	DnsCache = obliviousmap.NewObliviousMap[*interfaces.IPStacks]("DnsCache/", time.Minute, true, &memIPStacks)
}
