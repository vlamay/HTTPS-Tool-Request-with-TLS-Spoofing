package request

import (
	"bufio"
	"context"
	"custom-https-tool/internal/proxy"
	"custom-https-tool/internal/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

type Requestor struct {
	Client *http.Client
}

func NewRequestor(proxy *proxy.Proxy, tlsProfile string) (*Requestor, error) {
	helloID, err := tls.GetClientHelloSpec(tlsProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS client hello spec: %w", err)
	}

	dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		log.Printf("DEBUG: Dialing proxy %s for target %s", proxy.URL.Host, addr)
		proxyConn, err := net.DialTimeout("tcp", proxy.URL.Host, 60*time.Second)
		if err != nil {
			return nil, fmt.Errorf("failed to dial proxy: %w", err)
		}

		connectReq := &http.Request{
			Method: "CONNECT",
			URL:    &url.URL{Host: addr},
			Host:   addr,
			Header: make(http.Header),
		}
		connectReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
		connectReq.Header.Set("Proxy-Connection", "Keep-Alive")
		if proxy.URL.User != nil {
			if password, ok := proxy.URL.User.Password(); ok {
				auth := base64.StdEncoding.EncodeToString([]byte(proxy.URL.User.Username() + ":" + password))
				connectReq.Header.Set("Proxy-Authorization", "Basic "+auth)
			}
		}

		log.Printf("DEBUG: Sending CONNECT request:\n%s", requestToString(connectReq))

		if err := connectReq.Write(proxyConn); err != nil {
			return nil, fmt.Errorf("failed to write CONNECT request: %w", err)
		}

		br := bufio.NewReader(proxyConn)
		resp, err := http.ReadResponse(br, connectReq)
		if err != nil {
			return nil, fmt.Errorf("failed to read CONNECT response: %w", err)
		}
		defer resp.Body.Close()

		log.Printf("DEBUG: Received CONNECT response:\n%s", responseToString(resp))

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("proxy CONNECT request failed: %s", resp.Status)
		}

		uconn := utls.UClient(proxyConn, &utls.Config{ServerName: strings.Split(addr, ":")[0]}, helloID)
		if err := uconn.Handshake(); err != nil {
			return nil, err
		}
		return uconn, nil
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	transport := &http.Transport{
		DialContext: dialer,
	}

	if err := http2.ConfigureTransport(transport); err != nil {
		return nil, fmt.Errorf("failed to configure HTTP/2 transport: %w", err)
	}

	return &Requestor{
		Client: &http.Client{
			Transport: transport,
			Jar:       jar,
			Timeout:   60 * time.Second,
		},
	}, nil
}

func (r *Requestor) Do(req *http.Request) (*http.Response, error) {
	return r.Client.Do(req)
}

func requestToString(req *http.Request) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s HTTP/1.1\r\n", req.Method, req.URL.Host))
	sb.WriteString(fmt.Sprintf("Host: %s\r\n", req.Host))
	for k, v := range req.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, ", ")))
	}
	sb.WriteString("\r\n")
	return sb.String()
}

func responseToString(resp *http.Response) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s %s\r\n", resp.Proto, resp.Status))
	for k, v := range resp.Header {
		sb.WriteString(fmt.Sprintf("%s: %s\r\n", k, strings.Join(v, ", ")))
	}
	sb.WriteString("\r\n")
	return sb.String()
}
