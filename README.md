# Custom HTTPS Request Tool

This tool is designed to make custom HTTPS requests with TLS fingerprint spoofing to bypass anti-bot protections.

## Features

- Customizable HTTPS requests
- TLS fingerprint spoofing (Chrome, Firefox, Safari, random)
- Proxy support (HTTP, SOCKS5) with rotation
- User-friendly interface
- High reliability and speed

## Build and Run

1. **Install Go:** Make sure you have Go 1.20+ installed.
2. **Install dependencies:**
   ```bash
   go mod tidy
   ```
3. **Run the tool:**
   ```bash
   go run cmd/main.go
   ```

## Configuration

The configuration is loaded from the `config.json` file. Here's an example:

```json
{
  "target_urls": [
    "https://www.viagogo.co.uk/Concert-Tickets/Rock-and-Pop/Sting-Tickets/E-157332132",
    "https://www.viagogo.co.uk/Concert-Tickets/Alternative-and-Indie/Coldplay-Tickets/E-155391504",
    "https://www.stubhub.com/stardew-valley-denver-tickets-9-13-2025/event/156264784",
    "https://www.viagogo.com/Concert-Tickets/Alternative-Music/Coldplay-Tickets/E-155741198"
  ],
  "proxy_list": [
    "user:pass@ip:port"
  ],
  "num_requests": 50,
  "tls_profile": "random",
  "delay_range": [
    500,
    3000
  ]
}
```

- `target_urls`: A list of URLs to target.
- `proxy_list`: A list of proxies in the format `user:pass@ip:port`.
- `num_requests`: The total number of requests to make.
- `tls_profile`: The TLS profile to use. Can be "chrome", "firefox", "safari", or "random".
- `delay_range`: The minimum and maximum delay between requests in milliseconds.
