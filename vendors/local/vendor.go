package local

import (
	"context"
	"fmt"
	"net"

	"github.com/miaokobot/miaospeed/interfaces"
)

type Local struct {
	name string
	info string
}

func (c *Local) Type() interfaces.VendorType {
	return interfaces.VendorLocal
}

func (c *Local) Status() interfaces.VendorStatus {
	return interfaces.VStatusOperational
}

func (c *Local) Build(proxyName string, proxyInfo string) interfaces.Vendor {
	c.name = proxyName
	c.info = proxyInfo
	return c
}

func (c *Local) DialTCP(ctx context.Context, url string, network interfaces.RequestOptionsNetwork) (net.Conn, error) {
	if hostname, port, err := urlToMetadata(url); err != nil {
		return nil, err
	} else {
		return net.Dial(network.String(), fmt.Sprintf("%s:%d", hostname, port))
	}
}

func (c *Local) DialUDP(ctx context.Context, url string) (net.PacketConn, error) {
	return nil, fmt.Errorf("local test does not support udp yet")

}
func (c *Local) ProxyInfo() interfaces.ProxyInfo {
	return interfaces.ProxyInfo{
		Name:    c.name,
		Address: "127.0.0.1",
		Type:    interfaces.Parse("Invalid"),
	}
}
