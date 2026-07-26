package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"sync"
	"time"

	adapterhttp "github.com/Unikyri/vivi-perfilamiento-leads/internal/adapters/http"
)

type localServer struct {
	mu    sync.Mutex
	turns map[string]int
}

func newHandler() http.Handler {
	h := &localServer{turns: make(map[string]int)}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /salud", h.health)
	mux.HandleFunc("GET /api/leads", h.leads)
	mux.HandleFunc("POST /api/conversations", h.conversation)
	trusted := []netip.Prefix{netip.MustParsePrefix("127.0.0.1/32"), netip.MustParsePrefix("::1/128")}
	return adapterhttp.NuevoLimitadorTasa(mux, trusted)
}

func (s *localServer) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"estado": "OK", "bd": "in-memory", "llm": "stub", "red": "disabled",
	})
}

func (s *localServer) leads(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"total": 3,
		"leads": []map[string]string{
			{"id": "stub-ana", "estado": "PERFILANDO"},
			{"id": "stub-carlos", "estado": "PERFILANDO"},
			{"id": "stub-luisa", "estado": "PERFILANDO"},
		},
	})
}

func (s *localServer) conversation(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ClientID string `json:"client_id"`
		Message  string `json:"message"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.ClientID == "" || request.Message == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "client_id and message are required"})
		return
	}
	s.mu.Lock()
	s.turns[request.ClientID]++
	turn := s.turns[request.ClientID]
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{
		"client_id": request.ClientID,
		"turn":      turn,
		"reply":     "Respuesta determinística del LLM stub.",
		"provider":  "stub",
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func main() {
	port := flag.Int("port", 8080, "loopback TCP port")
	flag.Parse()
	if *port < 1 || *port > 65535 {
		log.Fatal("port must be between 1 and 65535")
	}
	server := &http.Server{
		Addr:              net.JoinHostPort("127.0.0.1", strconv.Itoa(*port)),
		Handler:           newHandler(),
		ReadHeaderTimeout: 2 * time.Second,
	}
	log.Printf("local load harness listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
