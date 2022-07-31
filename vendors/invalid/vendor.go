package invalid

import (
	"context"
	"fmt"
	"net"

	"github.com/miaokobot/miaospeed/interfaces"
)

type Invalid struct {
	name string
}

func (c *Invalid) Type() interfaces.VendorType {
	return interfaces.VendorInvalid
}

func (c *Invalid) Status() interfaces.VendorStatus {
	return interfaces.VStatusNotReady
}

func (c *Invalid) Build(proxyName string, proxyInfo string) interfaces.Vendor {
	c.name = proxyName
	return c
}

func (c *Invalid) DialTCP(ctx context.Context, url string, network interfaces.RequestOptionsNetwork) (net.Conn, error) {
	return nil, fmt.Errorf("the vendor is invalid")
}

func (c *Invalid) DialUDP(ctx context.Context, url string) (net.PacketConn, error) {
	return nil, fmt.Errorf("the vendor is invalid")
}

func (c *Invalid) ProxyInfo() interfaces.ProxyInfo {
	return interfaces.ProxyInfo{
		Name: c.name,
		Type: interfaces.ProxyInvalid,
	}
}
