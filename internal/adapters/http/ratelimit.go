package http

import (
	"errors"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"
)

const (
	defaultRateLimit  = 30
	defaultRateWindow = time.Minute
	defaultMaxClients = 4096
)

var ErrLimiteTasaHTTP = errors.New("http rate limit exceeded")

type rateWindow struct {
	started time.Time
	count   int
}
type rateLimiter struct {
	next       http.Handler
	now        func() time.Time
	limit      int
	window     time.Duration
	maxClients int
	trusted    []netip.Prefix
	mu         sync.Mutex
	clients    map[string]rateWindow
}

// NuevoLimitadorTasa applies the safe process-local API policy at the outer
// handler boundary. An empty trusted list ignores all forwarding headers.
func NuevoLimitadorTasa(next http.Handler, trusted []netip.Prefix) http.Handler {
	return nuevoLimitadorTasa(next, trusted, time.Now, defaultRateLimit, defaultRateWindow, defaultMaxClients)
}

func NuevoLimitadorTasaConLimite(next http.Handler, trusted []netip.Prefix, limit int) http.Handler {
	return nuevoLimitadorTasa(next, trusted, time.Now, limit, defaultRateWindow, defaultMaxClients)
}
func nuevoLimitadorTasa(next http.Handler, trusted []netip.Prefix, now func() time.Time, limit int, window time.Duration, maxClients int) http.Handler {
	if now == nil {
		now = time.Now
	}
	if limit < 1 {
		limit = defaultRateLimit
	}
	if window <= 0 {
		window = defaultRateWindow
	}
	if maxClients < 1 {
		maxClients = defaultMaxClients
	}
	return &rateLimiter{next: next, now: now, limit: limit, window: window, maxClients: maxClients, trusted: trusted, clients: make(map[string]rateWindow)}
}
func (l *rateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !esRutaAPI(r.URL.Path) || l.allow(identidad(r, l.trusted), l.now()) {
		l.next.ServeHTTP(w, r)
		return
	}
	writeError(w, ErrLimiteTasaHTTP)
}
func (l *rateLimiter) allow(client string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	for key, entry := range l.clients {
		if !now.Before(entry.started.Add(l.window)) {
			delete(l.clients, key)
		}
	}
	entry, exists := l.clients[client]
	if !exists {
		if len(l.clients) >= l.maxClients {
			return false
		}
		l.clients[client] = rateWindow{started: now, count: 1}
		return true
	}
	if entry.count >= l.limit {
		return false
	}
	entry.count++
	l.clients[client] = entry
	return true
}
func esRutaAPI(path string) bool { return path == "/api" || strings.HasPrefix(path, "/api/") }
func identidad(r *http.Request, trusted []netip.Prefix) string {
	peer, ok := parseAddr(r.RemoteAddr)
	if !ok {
		return strings.TrimSpace(r.RemoteAddr)
	}
	peer = peer.Unmap()
	if !isTrusted(peer, trusted) {
		return peer.String()
	}
	if value := r.Header.Get("Forwarded"); value != "" {
		if addresses, valid := parseForwarded(value); valid {
			return selectedAddress(addresses, trusted, peer)
		}
		return peer.String()
	}
	if value := r.Header.Get("X-Forwarded-For"); value != "" {
		if addresses, valid := parseXForwardedFor(value); valid {
			return selectedAddress(addresses, trusted, peer)
		}
	}
	return peer.String()
}
func selectedAddress(addresses []netip.Addr, trusted []netip.Prefix, peer netip.Addr) string {
	for i := len(addresses) - 1; i >= 0; i-- {
		address := addresses[i].Unmap()
		if !isTrusted(address, trusted) {
			return address.String()
		}
	}
	if len(addresses) > 0 {
		return addresses[0].Unmap().String()
	}
	return peer.String()
}
func parseAddr(value string) (netip.Addr, bool) {
	value = strings.TrimSpace(value)
	if host, _, err := net.SplitHostPort(value); err == nil {
		value = host
	} else {
		value = strings.Trim(value, "[]")
	}
	address, err := netip.ParseAddr(value)
	return address, err == nil
}
func parseXForwardedFor(value string) ([]netip.Addr, bool) {
	parts := strings.Split(value, ",")
	addresses := make([]netip.Addr, 0, len(parts))
	for _, part := range parts {
		address, ok := parseAddr(strings.TrimSpace(part))
		if !ok {
			return nil, false
		}
		addresses = append(addresses, address)
	}
	return addresses, len(addresses) > 0
}
func parseForwarded(value string) ([]netip.Addr, bool) {
	parts := strings.Split(value, ",")
	addresses := make([]netip.Addr, 0, len(parts))
	for _, part := range parts {
		var found string
		for _, parameter := range strings.Split(part, ";") {
			keyValue := strings.SplitN(strings.TrimSpace(parameter), "=", 2)
			if len(keyValue) != 2 || !strings.EqualFold(keyValue[0], "for") || found != "" {
				if len(keyValue) == 2 && strings.EqualFold(keyValue[0], "for") {
					return nil, false
				}
				continue
			}
			found = strings.Trim(strings.TrimSpace(keyValue[1]), "\"")
		}
		if found == "" {
			return nil, false
		}
		address, ok := parseAddr(found)
		if !ok {
			return nil, false
		}
		addresses = append(addresses, address)
	}
	return addresses, len(addresses) > 0
}
func isTrusted(address netip.Addr, trusted []netip.Prefix) bool {
	for _, prefix := range trusted {
		if prefix.Contains(address) || prefix.Contains(address.Unmap()) {
			return true
		}
	}
	return false
}
