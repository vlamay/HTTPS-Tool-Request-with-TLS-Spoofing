package main

import (
	"custom-https-tool/internal/config"
	"custom-https-tool/internal/proxy"
	"custom-https-tool/internal/request"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	proxyRotator, err := proxy.NewRotator(cfg.ProxyList)
	if err != nil {
		log.Fatalf("Failed to create proxy rotator: %v", err)
	}

	for i := 0; i < cfg.NumRequests; i++ {
		targetURL := cfg.TargetURLs[rand.Intn(len(cfg.TargetURLs))]
		proxy := proxyRotator.GetProxy()
		if proxy == nil {
			log.Println("No proxies available, stopping.")
			break
		}

		log.Printf("Request %d: URL=%s, Proxy=%s, TLSProfile=%s\n", i+1, targetURL, proxy.URL.String(), cfg.TLSProfile)

		requestor, err := request.NewRequestor(proxy, cfg.TLSProfile)
		if err != nil {
			log.Printf("Failed to create requestor: %v", err)
			continue
		}

		req, err := http.NewRequest("GET", targetURL, nil)
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			continue
		}

		// Mimic browser headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
		req.Header.Set("Sec-Fetch-Dest", "document")
		req.Header.Set("Sec-Fetch-Mode", "navigate")
		req.Header.Set("Sec-Fetch-Site", "none")
		req.Header.Set("Sec-Fetch-User", "?1")
		req.Header.Set("Upgrade-Insecure-Requests", "1")

		resp, err := requestor.Do(req)
		if err != nil {
			log.Printf("Request failed: %v", err)
			continue
		}
		defer resp.Body.Close()

		log.Printf("Response status: %s", resp.Status)

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			continue
		}

		// Challenge detection
		if resp.StatusCode == 403 || strings.Contains(string(body), "cf-challenge") || strings.Contains(string(body), "g-recaptcha") {
			log.Println("Challenge detected, retrying with new proxy and fingerprint.")
			continue
		}

		log.Println("Request successful.")

		// Random delay
		delay := time.Duration(rand.Intn(cfg.DelayRange[1]-cfg.DelayRange[0]+1)+cfg.DelayRange[0]) * time.Millisecond
		time.Sleep(delay)
	}
}
