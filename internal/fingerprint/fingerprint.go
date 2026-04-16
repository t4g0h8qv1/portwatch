// Package fingerprint identifies services running on open ports
// by reading banners or matching known port-to-service mappings.
package fingerprint

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// Result holds the fingerprint information for a single port.
type Result struct {
	Port    int    `json:"port"`
	Service string `json:"service"`
	Banner  string `json:"banner,omitempty"`
}

// wellKnown maps common ports to service names.
var wellKnown = map[int]string{
	21:   "ftp",
	22:   "ssh",
	23:   "telnet",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgresql",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Identify attempts to fingerprint the service on the given host:port.
// It first checks the well-known map, then tries to grab a banner.
func Identify(host string, port int, timeout time.Duration) Result {
	res := Result{Port: port}

	if svc, ok := wellKnown[port]; ok {
		res.Service = svc
	} else {
		res.Service = "unknown"
	}

	banner := grabBanner(host, port, timeout)
	if banner != "" {
		res.Banner = banner
		if res.Service == "unknown" {
			res.Service = inferFromBanner(banner)
		}
	}

	return res
}

func grabBanner(host string, port int, timeout time.Duration) string {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return ""
	}
	defer conn.Close()
	_ = conn.SetReadDeadline(time.Now().Add(timeout))

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil || n == 0 {
		return ""
	}
	return strings.TrimSpace(string(buf[:n]))
}

func inferFromBanner(banner string) string {
	lower := strings.ToLower(banner)
	switch {
	case strings.Contains(lower, "ssh"):
		return "ssh"
	case strings.Contains(lower, "ftp"):
		return "ftp"
	case strings.Contains(lower, "smtp"), strings.Contains(lower, "220"):
		return "smtp"
	case strings.Contains(lower, "http"):
		return "http"
	}
	return "unknown"
}
