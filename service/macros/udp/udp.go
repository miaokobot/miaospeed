package udp

import (
	"context"
	"strings"
	"sync"

	"github.com/miaokobot/miaospeed/interfaces"
)

func detectNATType(proxy interfaces.Vendor, url string) (nmt NATMapType, nft NATFilterType) {
	addrStr := strings.TrimLeft(url, "udp://")

	wg := sync.WaitGroup{}
	ctx := context.Background()

	wg.Add(1)
	go func() {
		if instance, _ := proxy.DialUDP(ctx, url); instance != nil {
			nmt = MappingTests(instance, addrStr)
			instance.Close()
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if instance, _ := proxy.DialUDP(ctx, url); instance != nil {
			nft = FilteringTests(instance, addrStr)
			instance.Close()
		}
		wg.Done()
	}()

	wg.Wait()
	return
}

func natTypeToString(nmt NATMapType, nft NATFilterType) string {
	if nmt == NATMapFailed || nft == NATFilterFailed {
		return "Unknown"
	}

	if nmt == NATMapIndependent {
		if nft == NATFilterIndependent {
			return "FullCone"
		} else if nft == NATFilterAddrIndependent {
			return "RestrictedCone"
		} else {
			return "PortRestrictedCone"
		}
	}

	if nmt == NATMapAddrPortIndependent && nft == NATFilterAddrPortIndependent {
		return "Symmetric"
	}
	return "SymmetricFirewall"
}
