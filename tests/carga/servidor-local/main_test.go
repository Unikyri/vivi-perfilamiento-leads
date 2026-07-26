package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func request(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.RemoteAddr = "127.0.0.1:43210"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestLocalRoutesReturnDeterministicResponses(t *testing.T) {
	h := newHandler()
	cases := []struct {
		name, method, path, body, want string
	}{
		{"health", http.MethodGet, "/salud", "", `"red":"disabled"`},
		{"leads", http.MethodGet, "/api/leads", "", `"total":3`},
		{"conversation", http.MethodPost, "/api/conversations", `{"client_id":"client-1","message":"hola"}`, `"provider":"stub"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response := request(h, tc.method, tc.path, tc.body)
			if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), tc.want) {
				t.Fatalf("status=%d body=%s", response.Code, response.Body)
			}
		})
	}
}

func TestLocalRoutesUseProductionRateLimitPolicy(t *testing.T) {
	h := newHandler()
	for i := 0; i < 30; i++ {
		if response := request(h, http.MethodGet, "/api/leads", ""); response.Code != http.StatusOK {
			t.Fatalf("request %d status=%d", i+1, response.Code)
		}
	}
	response := request(h, http.MethodGet, "/api/leads", "")
	if response.Code != http.StatusTooManyRequests || !strings.Contains(response.Body.String(), `"codigo":"LIMITE_TASA"`) {
		t.Fatalf("status=%d body=%s", response.Code, response.Body)
	}
	if health := request(h, http.MethodGet, "/salud", ""); health.Code != http.StatusOK {
		t.Fatalf("health status=%d", health.Code)
	}
}

func TestConversationRequiresPayload(t *testing.T) {
	response := request(newHandler(), http.MethodPost, "/api/conversations", `{}`)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", response.Code, response.Body)
	}
}
