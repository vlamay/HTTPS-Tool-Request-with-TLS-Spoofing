package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	TargetURLs  []string `json:"target_urls"`
	ProxyList   []string `json:"proxy_list"`
	NumRequests int      `json:"num_requests"`
	TLSProfile  string   `json:"tls_profile"`
	DelayRange  []int    `json:"delay_range"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
