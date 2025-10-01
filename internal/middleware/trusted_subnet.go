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
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if subnetIP == "" {
				http.Error(w, errForbidden.Error(), http.StatusForbidden)
				return
			}

			_, subnet, err := net.ParseCIDR(subnetIP)
			if err != nil {
				err := errParse.Error() + ": " + err.Error()
				http.Error(w, err, http.StatusInternalServerError)
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
