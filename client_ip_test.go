package httputil_test

import (
	"net/http"
	"testing"

	"github.com/d5/go-httputil"
)

func TestClientIPFromRequest(t *testing.T) {
	// empty
	testClientIPFromRequest(t, "", "", "")

	// no X-Forwarded-For
	testClientIPFromRequest(t, "29.14.99.122", "", "") // Request.RemoteAddr always have ":port"
	testClientIPFromRequest(t, "29.14.99.122:1984", "", "29.14.99.122")
	testClientIPFromRequest(t, "[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1123", "", "2001:db8:85a3:8d3:1319:8a2e:370:7348") // IPv6

	// X-Forwarded-For
	testClientIPFromRequest(t, "29.14.99.122:1984", "163.187.237.174", "163.187.237.174")
	testClientIPFromRequest(t, "29.14.99.122:1984", "163.187.237.174, 149.230.69.204", "149.230.69.204")
	testClientIPFromRequest(t, "29.14.99.122:1984", "163.187.237.174, 149.230.69.204, 102.10.86.163", "102.10.86.163")

	// X-Forwarded-For w/ private IPs
	testClientIPFromRequest(t, "29.14.99.122:1984", "163.187.237.174, 10.0.0.28", "163.187.237.174")
	testClientIPFromRequest(t, "29.14.99.122:1984", "163.187.237.174, 10.0.0.28, 192.168.0.4", "163.187.237.174")
}

func testClientIPFromRequest(t *testing.T, remoteAddr, forwardedFor, expected string) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com/", nil)
	req.RemoteAddr = remoteAddr
	req.Header.Set("X-Forwarded-For", forwardedFor)
	actual := httputil.ClientIPFromRequest(req)
	if actual != expected {
		t.Errorf("Expected: %q, Actual: %q\n", expected, actual)
	}
}
