package clash

import (
	"context"
	"fmt"
	"net"

	"github.com/Dreamacro/clash/constant"
	"github.com/miaokobot/miaospeed/interfaces"
)

type Clash struct {
	proxy constant.Proxy
}

func (c *Clash) Type() interfaces.VendorType {
	return interfaces.VenderClash
}

func (c *Clash) Status() interfaces.VendorStatus {
	if c == nil || c.proxy == nil {
		return interfaces.VStatusNotReady
	}

	return interfaces.VStatusOperational
}

func (c *Clash) Build(proxyName string, proxyInfo string) interfaces.Vendor {
	if c == nil {
		c = &Clash{}
	}
	c.proxy = extractFirstProxy(proxyName, proxyInfo)
	return c
}

func (c *Clash) DialTCP(ctx context.Context, url string, network interfaces.RequestOptionsNetwork) (net.Conn, error) {
	if c == nil || c.proxy == nil {
		return nil, fmt.Errorf("should call Build() before run")
	}

	addr, err := urlToMetadata(url, constant.TCP)
	if err != nil {
		return nil, fmt.Errorf("cannot build tcp context")
	}

	return c.proxy.DialContext(ctx, &addr)
}

func (c *Clash) DialUDP(ctx context.Context, url string) (net.PacketConn, error) {
	if c == nil || c.proxy == nil {
		return nil, fmt.Errorf("should call Build() before run")
	}

	addr, err := urlToMetadata(url, constant.UDP)
	if err != nil {
		return nil, fmt.Errorf("cannot build udp context")
	}

	return c.proxy.DialUDP(&addr)

}
func (c *Clash) ProxyInfo() interfaces.ProxyInfo {
	if c == nil || c.proxy == nil {
		return interfaces.ProxyInfo{}
	}

	return interfaces.ProxyInfo{
		Name:    c.proxy.Name(),
		Address: c.proxy.Addr(),
		Type:    interfaces.Parse(c.proxy.Type().String()),
	}
}
