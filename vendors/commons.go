package vendors

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/miaokobot/miaospeed/interfaces"
	"github.com/miaokobot/miaospeed/utils/structs"
)

// for all methods in commons
// if proxy is nil, fallback to system
// if proxy is not nil but it is not ready
// errors will be returned

func RequestUnsafe(ctx context.Context, p interfaces.Vendor, reqOpt *interfaces.RequestOptions) (*http.Response, []string, error) {
	if p != nil && p.Status() == interfaces.VStatusNotReady || reqOpt == nil {
		return nil, nil, errors.New("proxy is not ready")
	}

	// check request method
	if reqOpt.Method == "" {
		reqOpt.Method = http.MethodGet
	}

	// check body reader
	var reader io.Reader = nil
	if len(reqOpt.Body) > 0 {
		reader = bytes.NewBuffer(reqOpt.Body)
	}

	// build request
	req, err := http.NewRequest(reqOpt.Method, reqOpt.URL, reader)
	if err != nil {
		return nil, nil, err
	}

	// write headers and cookies
	for hkey, hval := range reqOpt.Headers {
		req.Header.Add(hkey, hval)
	}
	for ckey, cval := range reqOpt.Cookies {
		req.AddCookie(&http.Cookie{Name: ckey, Value: cval})
	}
	req = req.WithContext(ctx)

	// connect proxy bridge
	// init params copied from http.DefaultTransport
	transport := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if p != nil {
		transport.Dial = func(string, string) (net.Conn, error) {
			return p.DialTCP(ctx, reqOpt.URL, reqOpt.Network)
		}
	} else {
		transport.Dial = func(string, string) (net.Conn, error) {
			return net.Dial(reqOpt.Network.String(), reqOpt.URL)
		}
	}

	// make a list to record all redirects
	// the maximum count of redirection is 64
	redirects := []string{}
	client := http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if reqOpt.NoRedir || len(redirects) > 64 {
				return http.ErrUseLastResponse
			}

			redirects = append(redirects, req.Response.Header.Get("Location"))
			return nil
		},
	}

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return resp, redirects, nil
}

func Request(ctx context.Context, p interfaces.Vendor, reqOpt *interfaces.RequestOptions) (duration uint16, bodyBytes []byte, resp *http.Response, redirects []string) {
	var err error

	start := time.Now()
	resp, redirects, err = RequestUnsafe(ctx, p, reqOpt)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	duration = uint16(time.Since(start) / time.Millisecond)
	return
}

func RequestWithRetry(p interfaces.Vendor, retry int, timeoutMillisecond int64, reqOpt *interfaces.RequestOptions) ([]byte, *http.Response, []string) {
	var resp *http.Response = nil
	var retBody []byte = nil
	var redirects = []string{}

	for i := 0; resp == nil && i < structs.WithIn(retry, 1, 10); i += 1 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMillisecond)*time.Millisecond)
		_, retBody, resp, redirects = Request(ctx, p, reqOpt)
		cancel()
	}

	return retBody, resp, redirects
}

func NetCat(ctx context.Context, p interfaces.Vendor, addr string, data []byte, network interfaces.RequestOptionsNetwork) ([]byte, error) {
	var conn net.Conn = nil
	err := fmt.Errorf("proxy is not ready")
	if p != nil {
		if p.Status() == interfaces.VStatusOperational {
			conn, err = p.DialTCP(ctx, addr, network)
		}
	} else {
		conn, err = net.Dial(network.String(), addr)
	}

	if err != nil {
		return nil, err
	}

	defer conn.Close()
	if _, err := conn.Write(data); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, conn); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func NetCatWithRetry(p interfaces.Vendor, retry int, timeoutMillisecond int64, addr string, data []byte, network interfaces.RequestOptionsNetwork) ([]byte, error) {
	var retBody []byte = nil
	var err = fmt.Errorf("request not send")

	for i := 0; err != nil && i < structs.WithIn(retry, 1, 10); i += 1 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMillisecond)*time.Millisecond)
		retBody, err = NetCat(ctx, p, addr, data, network)
		cancel()
	}

	return retBody, err
}
