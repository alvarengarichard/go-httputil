package httputil

import (
	"net"
	"net/http"
	"strings"
)

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

func isPrivateIP(ip net.IP) bool {
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

// ClientIPFromRequest returns the public IP address of the client.
func ClientIPFromRequest(req *http.Request) string {
	addrs := strings.Split(req.Header.Get("X-Forwarded-For"), ",")

	// march from right to left until we get a public address
	// that will be the address right before our proxy.
	for i := len(addrs) - 1; i >= 0; i-- {
		ip := strings.TrimSpace(addrs[i])

		// header can contain spaces too, strip those out.
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil || !parsedIP.IsGlobalUnicast() || isPrivateIP(parsedIP) {
			continue
		}

		return ip
	}

	// use RemoteAddr ("IPv4:port" or "[IPv6]:port")
	if req.RemoteAddr != "" {
		host, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return ""
		}

		parsedIP := net.ParseIP(host)
		if parsedIP != nil {
			return parsedIP.String()
		}

		return host
	}

	return ""
}
