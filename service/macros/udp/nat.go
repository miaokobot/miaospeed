package udp

import (
	"errors"
	"net"
	"time"

	"github.com/miaokobot/miaospeed/utils"
	"github.com/pion/stun"
)

type NATMapType int

const (
	NATMapFailed NATMapType = iota
	NATMapIndependent
	NATMapAddrIndependent
	NATMapAddrPortIndependent
	NATMapNoNat
)

type NATFilterType int

const (
	NATFilterFailed NATFilterType = iota
	NATFilterIndependent
	NATFilterAddrIndependent
	NATFilterAddrPortIndependent
)

type stunServerConn struct {
	conn        net.PacketConn
	LocalAddr   net.Addr
	RemoteAddr  *net.UDPAddr
	OtherAddr   *net.UDPAddr
	messageChan chan *stun.Message
}

func (c *stunServerConn) Close() error {
	return nil
}

const (
	messageHeaderSize = 20
	natTimeout        = 3
)

var (
	errResponseMessage = errors.New("error reading from response message channel")
	errTimedOut        = errors.New("timed out waiting for response")
	errNoOtherAddress  = errors.New("no OTHER-ADDRESS in message")
)

// RFC5780: 4.3.  Determining NAT Mapping Behavior
func MappingTests(conn net.PacketConn, addrStr string) NATMapType {
	mapTestConn, err := connect(conn, addrStr)
	if err != nil {
		utils.DLog("NAT MAP TEST | cannot connect to stun server:", err.Error())
		return NATMapFailed
	}

	// Test I: Regular binding request
	request := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	resp, err := mapTestConn.roundTrip(request, mapTestConn.RemoteAddr, natTimeout)
	if err != nil {
		utils.DLog("NAT MAP TEST | TEST I Failed:", err.Error())
		return NATMapFailed
	}

	// Parse response message for XOR-MAPPED-ADDRESS and make sure OTHER-ADDRESS valid
	resps1 := parse(resp)
	if resps1.xorAddr == nil || resps1.otherAddr == nil {
		utils.DLog("NAT MAP TEST | TEST I Failed: no other address")
		return NATMapFailed
	}
	addr, err := net.ResolveUDPAddr("udp4", resps1.otherAddr.String())
	if err != nil {
		utils.DLog("NAT MAP TEST | TEST I Resolve Failed:", err.Error())
		return NATMapFailed
	}
	mapTestConn.OtherAddr = addr

	// Assert mapping behavior
	if resps1.xorAddr.String() == mapTestConn.LocalAddr.String() {
		return NATMapNoNat
	}

	// Test II: Send binding request to the other address but primary port
	oaddr := *mapTestConn.OtherAddr
	oaddr.Port = mapTestConn.RemoteAddr.Port
	resp, err = mapTestConn.roundTrip(request, &oaddr, natTimeout)
	if err != nil {
		utils.DLog("NAT MAP TEST | TEST II Failed:", err.Error())
		return NATMapFailed
	}

	// Assert mapping behavior
	resps2 := parse(resp)
	if resps2.xorAddr.String() == resps1.xorAddr.String() {
		return NATMapIndependent
	}

	// Test III: Send binding request to the other address and port
	resp, err = mapTestConn.roundTrip(request, mapTestConn.OtherAddr, natTimeout)
	if err != nil {
		utils.DLog("NAT MAP TEST | TEST III Failed:", err.Error())
		return NATMapFailed
	}

	// Assert mapping behavior
	resps3 := parse(resp)
	if resps3.xorAddr.String() == resps2.xorAddr.String() {
		return NATMapAddrIndependent
	} else {
		return NATMapAddrPortIndependent
	}
}

// RFC5780: 4.4.  Determining NAT Filtering Behavior
func FilteringTests(conn net.PacketConn, addrStr string) NATFilterType {
	mapTestConn, err := connect(conn, addrStr)
	if err != nil {
		utils.DLog("NAT FLT TEST | cannot connect to stun server:", err.Error())
		return NATFilterFailed
	}

	// Test I: Regular binding request
	request := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	resp, err := mapTestConn.roundTrip(request, mapTestConn.RemoteAddr, natTimeout)
	if err != nil || errors.Is(err, errTimedOut) {
		utils.DLog("NAT FLT TEST | TEST I Failed:", err.Error())
		return NATFilterFailed
	}
	resps := parse(resp)
	if resps.xorAddr == nil || resps.otherAddr == nil {
		utils.DLog("NAT FLT TEST | TEST I Failed: no other address")
		return NATFilterFailed
	}
	addr, err := net.ResolveUDPAddr("udp4", resps.otherAddr.String())
	if err != nil {
		utils.DLog("NAT FLT TEST | TEST I Failed:", err.Error())
		return NATFilterFailed
	}
	mapTestConn.OtherAddr = addr

	// Test II: Request to change both IP and port
	request = stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	request.Add(stun.AttrChangeRequest, []byte{0x00, 0x00, 0x00, 0x06})

	resp, err = mapTestConn.roundTrip(request, mapTestConn.RemoteAddr, natTimeout)
	if err == nil {
		return NATFilterIndependent
	} else if !errors.Is(err, errTimedOut) {
		utils.DLog("NAT FLT TEST | TEST II Failed:", err.Error())
		return NATFilterFailed
	}

	// Test III: Request to change port only
	request = stun.MustBuild(stun.TransactionID, stun.BindingRequest)
	request.Add(stun.AttrChangeRequest, []byte{0x00, 0x00, 0x00, 0x02})
	resp, err = mapTestConn.roundTrip(request, mapTestConn.RemoteAddr, natTimeout)
	if err == nil {
		return NATFilterAddrIndependent
	} else if errors.Is(err, errTimedOut) {
		return NATFilterAddrPortIndependent
	}

	return NATFilterFailed
}

// Parse a STUN message
func parse(msg *stun.Message) (ret struct {
	xorAddr    *stun.XORMappedAddress
	otherAddr  *stun.OtherAddress
	respOrigin *stun.ResponseOrigin
	mappedAddr *stun.MappedAddress
	software   *stun.Software
}) {
	ret.mappedAddr = &stun.MappedAddress{}
	ret.xorAddr = &stun.XORMappedAddress{}
	ret.respOrigin = &stun.ResponseOrigin{}
	ret.otherAddr = &stun.OtherAddress{}
	ret.software = &stun.Software{}
	if ret.xorAddr.GetFrom(msg) != nil {
		ret.xorAddr = nil
	}
	if ret.otherAddr.GetFrom(msg) != nil {
		ret.otherAddr = nil
	}
	if ret.respOrigin.GetFrom(msg) != nil {
		ret.respOrigin = nil
	}
	if ret.mappedAddr.GetFrom(msg) != nil {
		ret.mappedAddr = nil
	}
	if ret.software.GetFrom(msg) != nil {
		ret.software = nil
	}
	return ret
}

// Given an address string, returns a StunServerConn
func connect(conn net.PacketConn, addrStr string) (*stunServerConn, error) {
	addr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		return nil, err
	}

	mChan := listen(conn)
	return &stunServerConn{
		conn:        conn,
		LocalAddr:   conn.LocalAddr(),
		RemoteAddr:  addr,
		messageChan: mChan,
	}, nil
}

// Send request and wait for response or timeout
func (c *stunServerConn) roundTrip(msg *stun.Message, addr net.Addr, timeout int) (*stun.Message, error) {
	_ = msg.NewTransactionID()
	_, err := c.conn.WriteTo(msg.Raw, addr)
	if err != nil {
		return nil, err
	}

	// Wait for response or timeout
	select {
	case m, ok := <-c.messageChan:
		if !ok {
			return nil, errResponseMessage
		}
		return m, nil
	case <-time.After(time.Duration(timeout) * time.Second):
		return nil, errTimedOut
	}
}

// taken from https://github.com/pion/stun/blob/master/cmd/stun-traversal/main.go
func listen(conn net.PacketConn) (messages chan *stun.Message) {
	messages = make(chan *stun.Message)
	go func() {
		for {
			buf := make([]byte, 1024)

			n, _, err := conn.ReadFrom(buf)
			if err != nil {
				close(messages)
				return
			}
			buf = buf[:n]

			m := new(stun.Message)
			m.Raw = buf
			err = m.Decode()
			if err != nil {
				close(messages)
				return
			}

			messages <- m
		}
	}()
	return
}
