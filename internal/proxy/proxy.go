package proxy

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type Proxy struct {
	URL *url.URL
}

func NewProxy(proxyURL string) (*Proxy, error) {
	if !strings.HasPrefix(proxyURL, "http://") && !strings.HasPrefix(proxyURL, "https://") && !strings.HasPrefix(proxyURL, "socks5://") {
		proxyURL = "http://" + proxyURL
	}
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
	}
	return &Proxy{URL: parsedURL}, nil
}

type Rotator struct {
	proxies []*Proxy
	r       *rand.Rand
}

func NewRotator(proxyURLs []string) (*Rotator, error) {
	var proxies []*Proxy
	for _, proxyURL := range proxyURLs {
		p, err := NewProxy(proxyURL)
		if err != nil {
			return nil, err
		}
		proxies = append(proxies, p)
	}
	return &Rotator{
		proxies: proxies,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (r *Rotator) GetProxy() *Proxy {
	if len(r.proxies) == 0 {
		return nil
	}
	return r.proxies[r.r.Intn(len(r.proxies))]
}
