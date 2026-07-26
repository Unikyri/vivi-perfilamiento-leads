package http

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fakeClock struct{ now time.Time }

func (c *fakeClock) Now() time.Time { return c.now }

func testLimiter(clock *fakeClock, max int, trusted ...string) (*rateLimiter, *atomic.Int32) {
	var calls atomic.Int32
	prefixes := make([]netip.Prefix, 0, len(trusted))
	for _, value := range trusted {
		prefixes = append(prefixes, netip.MustParsePrefix(value))
	}
	next := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { calls.Add(1) })
	handler := nuevoLimitadorTasa(next, prefixes, clock.Now, 30, time.Minute, max).(*rateLimiter)
	return handler, &calls
}

func request(remote, path string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, path, nil)
	r.RemoteAddr = remote
	return r
}

func TestRateLimiterRoutesAndWindow(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	limiter, calls := testLimiter(clock, 10)
	for i := 0; i < 30; i++ {
		recorder := httptest.NewRecorder()
		limiter.ServeHTTP(recorder, request("192.0.2.1:1234", "/api/leads"))
		if recorder.Code != http.StatusOK {
			t.Fatalf("request %d status=%d", i+1, recorder.Code)
		}
	}
	blocked := httptest.NewRecorder()
	limiter.ServeHTTP(blocked, request("192.0.2.1:1234", "/api/leads"))
	if blocked.Code != http.StatusTooManyRequests || !strings.Contains(blocked.Body.String(), `"codigo":"LIMITE_TASA"`) {
		t.Fatalf("blocked response: %d %s", blocked.Code, blocked.Body)
	}
	if strings.Contains(blocked.Body.String(), "192.0.2.1") || strings.Contains(blocked.Body.String(), "trusted") {
		t.Fatal("rejection leaked limiter state")
	}
	if calls.Load() != 30 {
		t.Fatalf("downstream calls=%d, want 30", calls.Load())
	}

	other := httptest.NewRecorder()
	limiter.ServeHTTP(other, request("192.0.2.2:1234", "/api/leads"))
	if other.Code != http.StatusOK {
		t.Fatalf("isolated identity status=%d", other.Code)
	}
	for _, path := range []string{"/salud", "/", "/assets/app.js"} {
		for i := 0; i < 40; i++ {
			limiter.ServeHTTP(httptest.NewRecorder(), request("192.0.2.1:1234", path))
		}
	}
	clock.now = clock.now.Add(time.Minute)
	fresh := httptest.NewRecorder()
	limiter.ServeHTTP(fresh, request("192.0.2.1:1234", "/api/leads"))
	if fresh.Code != http.StatusOK {
		t.Fatalf("exact window reset status=%d", fresh.Code)
	}
}

func TestIdentityTrustAndCanonicalization(t *testing.T) {
	clock := &fakeClock{now: time.Now()}
	limiter, _ := testLimiter(clock, 100, "127.0.0.1/32", "10.0.0.0/8")
	untrusted := request("198.51.100.10:1", "/api")
	untrusted.Header.Set("X-Forwarded-For", "203.0.113.1")
	if got := identidad(untrusted, nil); got != "198.51.100.10" {
		t.Fatalf("untrusted identity=%q", got)
	}
	trusted := request("127.0.0.1:1", "/api")
	trusted.Header.Set("X-Forwarded-For", "198.51.100.1, 10.1.1.1")
	if got := identidad(trusted, limiter.trusted); got != "198.51.100.1" {
		t.Fatalf("trusted xff identity=%q", got)
	}
	trusted.Header.Set("Forwarded", `for=198.51.100.2;proto=https, for=10.1.1.1`)
	if got := identidad(trusted, limiter.trusted); got != "198.51.100.2" {
		t.Fatalf("trusted forwarded identity=%q", got)
	}
	trusted.Header.Set("Forwarded", "for=not-an-ip")
	if got := identidad(trusted, limiter.trusted); got != "127.0.0.1" {
		t.Fatalf("malformed identity=%q", got)
	}
	mapped := request("[::ffff:192.0.2.3]:1", "/api")
	if identidad(mapped, nil) != identidad(request("192.0.2.3:2", "/api"), nil) {
		t.Fatal("IPv4-mapped address was not canonicalized")
	}
}

func TestRateLimiterBoundedStateAndConcurrentAccess(t *testing.T) {
	clock := &fakeClock{now: time.Now()}
	limiter, _ := testLimiter(clock, 2)
	for _, peer := range []string{"192.0.2.1:1", "192.0.2.2:1"} {
		if !limiter.allow(peer, clock.now) {
			t.Fatal("initial identity rejected")
		}
	}
	if limiter.allow("192.0.2.3", clock.now) {
		t.Fatal("capacity overflow admitted")
	}
	clock.now = clock.now.Add(time.Minute)
	if !limiter.allow("192.0.2.3", clock.now) || len(limiter.clients) != 1 {
		t.Fatalf("expired state not reclaimed: len=%d", len(limiter.clients))
	}

	clock.now = time.Now()
	limiter, _ = testLimiter(clock, 4096)
	var allowed atomic.Int32
	var group sync.WaitGroup
	for i := 0; i < 100; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			if limiter.allow("192.0.2.9", clock.now) {
				allowed.Add(1)
			}
		}()
	}
	group.Wait()
	if allowed.Load() != 30 {
		t.Fatalf("concurrent allowed=%d, want 30", allowed.Load())
	}
}
