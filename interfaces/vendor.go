package interfaces

import (
	"context"
	"net"
)

type VendorType string

const (
	VendorLocal VendorType = "Local"
	VendorClash VendorType = "Clash"

	VendorInvalid VendorType = "Invalid"
)

type VendorStatus uint

const (
	VStatusOperational VendorStatus = iota

	VStatusNotReady
)

// a Vendor is an interface that allow macros to
// trigger connections from
type Vendor interface {
	// returns the type of the vendor
	Type() VendorType

	// returns the status of the vendor
	Status() VendorStatus

	// build conn based on proxy info string
	Build(proxyName string, proxyInfo string) Vendor

	// tcp connections
	DialTCP(ctx context.Context, url string, network RequestOptionsNetwork) (net.Conn, error)

	// udp connections
	DialUDP(ctx context.Context, url string) (net.PacketConn, error)

	// return universal proxy info
	ProxyInfo() ProxyInfo
}
