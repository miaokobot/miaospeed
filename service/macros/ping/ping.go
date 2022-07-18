package ping

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	urllib "net/url"
	"strings"
	"time"

	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/preconfigs"
	"github.com/miaokobot/miaospeed/utils"
	"github.com/miaokobot/miaospeed/utils/structs"
)

func pingViaTrace(ctx context.Context, p interfaces.Vendor, url string) (uint16, uint16, error) {
	transport := &http.Transport{
		Dial: func(string, string) (net.Conn, error) {
			return p.DialTCP(ctx, url, interfaces.ROptionsTCP)
		},
		// from http.DefaultTransport
		MaxIdleConns:          100,
		IdleConnTimeout:       3 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			// for version prior to tls1.3, the handshake will take 2-RTTs,
			// plus, majority server supports tls1.3, so we set a limit here
			MinVersion: tls.VersionTLS13,
			RootCAs:    preconfigs.MiaokoRootCAPrepare(),
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}

	tlsStart := int64(0)
	tlsEnd := int64(0)
	writeStart := int64(0)
	writeEnd := int64(0)
	trace := &httptrace.ClientTrace{
		TLSHandshakeStart: func() {
			tlsStart = time.Now().UnixMilli()
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			tlsEnd = time.Now().UnixMilli()
			if err != nil {
				tlsEnd = 0
			}
		},
		GotFirstResponseByte: func() {
			writeEnd = time.Now().UnixMilli()
		},
		WroteHeaders: func() {
			writeStart = time.Now().UnixMilli()
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	connStart := time.Now().UnixMilli()
	if resp, err := transport.RoundTrip(req); err != nil {
		return 0, 0, err
	} else {
		connEnd := time.Now().UnixMilli()
		utils.DBlackhole(!strings.HasPrefix(url, "https:"), connEnd-writeEnd, writeEnd-tlsEnd, tlsEnd-tlsStart, tlsStart-connStart)
		if !strings.HasPrefix(url, "https:") {
			return uint16(writeStart - connStart), uint16(writeEnd - connStart), nil
		}
		if resp.TLS != nil && resp.TLS.HandshakeComplete {
			// use payload rtt
			return uint16(writeEnd - tlsEnd), uint16(writeEnd - connStart), nil
			// return uint16(tlsEnd - tlsStart), uint16(writeEnd - connStart), nil
		}
		return 0, 0, fmt.Errorf("cannot extract payload from response")
	}
}

func pingViaNetCat(ctx context.Context, p interfaces.Vendor, url string) (uint16, uint16, error) {
	purl, _ := urllib.Parse(url)
	data := purl.EscapedPath()
	if purl.RawQuery != "" {
		data += "?" + purl.Query().Encode()
	}

	data = structs.X(preconfigs.NETCAT_HTTP_PAYLOAD, data, purl.Hostname(), utils.VERSION)

	connStart := time.Now().UnixMilli()
	conn, err := p.DialTCP(ctx, url, interfaces.ROptionsTCP)
	if err != nil || conn == nil {
		return 0, 0, fmt.Errorf("cannot dial remote address")
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second * 6))
	c := bufio.NewReader(conn)

	// prewrite to ensure tcp conn is established
	// httpStartReq1 := time.Now().UnixMilli()
	if _, err := conn.Write([]byte(data)); err != nil {
		return 0, 0, fmt.Errorf("cannot write payload to remote")
	}

	c.ReadLine()
	for c.Buffered() > 0 {
		c.ReadLine()
	}

	httpStartReq2 := time.Now().UnixMilli()
	if _, err := conn.Write([]byte(data)); err != nil {
		return 0, 0, fmt.Errorf("cannot write payload to remote")
	}

	c.ReadLine()
	for c.Buffered() > 0 {
		c.ReadLine()
	}
	httpEnd := time.Now().UnixMilli()

	return uint16(httpEnd - httpStartReq2), uint16(httpStartReq2 - connStart), nil
}

func ping(p interfaces.Vendor, url string, withAvg uint16, maxAttempt int, timeout uint) (uint16, uint16) {
	if p == nil {
		return 0, 0
	}

	failNum := 0
	totalMS := []uint16{}
	totalMSRTT := []uint16{}

	if withAvg < 1 || withAvg > uint16(maxAttempt) {
		withAvg = 1
	}

	for failNum+len(totalMS) < maxAttempt && len(totalMS) < int(withAvg) && maxAttempt-failNum >= int(withAvg) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
		delayRTT, delay := uint16(0), uint16(0)
		if strings.HasPrefix(url, "https:") {
			delayRTT, delay, _ = pingViaTrace(ctx, p, url)
		} else {
			delayRTT, delay, _ = pingViaNetCat(ctx, p, url)
		}
		if delayRTT > 0 {
			totalMSRTT = append(totalMSRTT, delayRTT)
			totalMS = append(totalMS, delay)
		} else {
			failNum += 1
		}
		cancel()
	}

	resultRTT, result := uint16(0), uint16(0)
	if len(totalMSRTT) >= int(withAvg) {
		resultRTT = computeAvgOfPing(totalMSRTT)
		result = computeAvgOfPing(totalMS)
	}
	return resultRTT, result
}
