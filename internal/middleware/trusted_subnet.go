package middleware

import (
	"errors"
	"net"
	"net/http"
)

var (
	errForbidden = errors.New("доступ запрещен")
	errParse     = errors.New("не удалось прочитать подсеть")
)

// CheckTrustedSubnet HTTP middleware that restricts access to handlers
// based on the client's IP address and a configured trusted subnet (CIDR).
func CheckTrustedSubnet(subnetIP string) func(http.Handler) http.Handler {
	var (
		subnet        *net.IPNet
		trustedSubnet bool
	)

	trustedSubnet = true

	if subnetIP == "" {
		trustedSubnet = false
	}

	_, subnet, err := net.ParseCIDR(subnetIP)
	if err != nil {
		trustedSubnet = false
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !trustedSubnet {
				http.Error(w, errForbidden.Error(), http.StatusForbidden)
				return
			}

			headerRealIP := r.Header.Get("X-Real-IP")
			if headerRealIP == "" {
				http.Error(w, errForbidden.Error(), http.StatusForbidden)
				return
			}

			ip := net.ParseIP(headerRealIP)
			if ip == nil || !subnet.Contains(ip) {
				http.Error(w, errForbidden.Error(), http.StatusForbidden)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
