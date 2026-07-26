package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type asyncProcessorFake struct {
	started chan struct{}
	release chan struct{}
	done    chan struct{}
	once    sync.Once
	mu      sync.Mutex
	input   usecase.EntradaMensaje
	err     error
}

type testID struct{ n int }

func (f *testID) Nuevo() string { f.n++; return "msg-" + strconv.Itoa(f.n) }

type testClock struct{ now time.Time }

func (f testClock) Ahora() time.Time         { return f.now }
func (f testClock) FechaSimulada() time.Time { return f.now }
func (f testClock) Avanzar(time.Time)        {}

func (f *asyncProcessorFake) Ejecutar(ctx context.Context, in usecase.EntradaMensaje) error {
	f.mu.Lock()
	f.input = in
	f.mu.Unlock()
	f.once.Do(func() { close(f.started) })
	select {
	case <-f.release:
		f.err = nil
	case <-ctx.Done():
		f.err = ctx.Err()
	}
	close(f.done)
	return f.err
}
func newAsyncProcessor() *asyncProcessorFake {
	return &asyncProcessorFake{started: make(chan struct{}), release: make(chan struct{}), done: make(chan struct{})}
}
func messageRouter(repo usecase.LeadRepository, tracker TurnoTracker) *nethttp.ServeMux {
	c, _ := NuevoControlador(Dependencias{Perfilar: &perfiladorHTTPFake{}, Leads: repo, Turnos: tracker})
	mux := nethttp.NewServeMux()
	c.Registrar(mux)
	return mux
}
func asyncLeadRepo() *leadHTTPFake {
	return &leadHTTPFake{lead: &domain.Lead{LeadID: "lead-1", Estado: domain.EstadoLeadPerfilando}}
}
func newTracker(p *asyncProcessorFake) *EjecutorTurnos {
	return NuevoEjecutorTurnos(p, &testID{}, testClock{now: time.Date(2026, 7, 26, 7, 0, 0, 0, time.UTC)})
}
func postMessage(mux *nethttp.ServeMux, body string) *httptest.ResponseRecorder {
	r := httptest.NewRecorder()
	mux.ServeHTTP(r, httptest.NewRequest(nethttp.MethodPost, "/api/leads/lead-1/mensajes", strings.NewReader(body)))
	return r
}
func waitInactive(t *testing.T, tracker *EjecutorTurnos) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for tracker.Activo("lead-1") {
		if time.Now().After(deadline) {
			t.Fatal("turn remained active")
		}
		time.Sleep(time.Millisecond)
	}
}

func TestS2MessageAcceptedAndPollsActiveToClear(t *testing.T) {
	p := newAsyncProcessor()
	tracker := newTracker(p)
	defer tracker.Cerrar()
	mux := messageRouter(asyncLeadRepo(), tracker)
	response := postMessage(mux, `{"tipo":"TEXTO","texto":"quiero comprar"}`)
	if response.Code != nethttp.StatusAccepted || !tracker.Activo("lead-1") {
		t.Fatalf("accept=%d active=%v", response.Code, tracker.Activo("lead-1"))
	}
	var accepted mensajeResponse
	if err := json.NewDecoder(response.Body).Decode(&accepted); err != nil || accepted.MensajeID != "msg-1" || !accepted.TurnoEnProceso || accepted.RecibidoEn.IsZero() {
		t.Fatalf("accepted=%+v err=%v", accepted, err)
	}
	<-p.started
	poll := httptest.NewRecorder()
	mux.ServeHTTP(poll, httptest.NewRequest(nethttp.MethodGet, "/api/leads/lead-1/conversacion", nil))
	var conversation conversacionResponse
	_ = json.NewDecoder(poll.Body).Decode(&conversation)
	if poll.Code != 200 || !conversation.TurnoEnProceso {
		t.Fatalf("poll active=%d %+v", poll.Code, conversation)
	}
	close(p.release)
	<-p.done
	waitInactive(t, tracker)
	poll = httptest.NewRecorder()
	mux.ServeHTTP(poll, httptest.NewRequest(nethttp.MethodGet, "/api/leads/lead-1/conversacion", nil))
	_ = json.NewDecoder(poll.Body).Decode(&conversation)
	if conversation.TurnoEnProceso {
		t.Fatal("poll remained active after processing")
	}
}

func TestS2RejectsConcurrentTurnAndCancelsOnClose(t *testing.T) {
	p := newAsyncProcessor()
	tracker := newTracker(p)
	mux := messageRouter(asyncLeadRepo(), tracker)
	if got := postMessage(mux, `{"tipo":"TEXTO","texto":"uno"}`).Code; got != 202 {
		t.Fatalf("first status=%d", got)
	}
	<-p.started
	second := postMessage(mux, `{"tipo":"TEXTO","texto":"dos"}`)
	if second.Code != 429 || !strings.Contains(second.Body.String(), `"codigo":"LIMITE_TASA"`) {
		t.Fatalf("second=%d %s", second.Code, second.Body)
	}
	tracker.Cerrar()
	select {
	case <-p.done:
	case <-time.After(time.Second):
		t.Fatal("processor did not cancel")
	}
	if !errors.Is(p.err, context.Canceled) || tracker.Activo("lead-1") {
		t.Fatalf("err=%v active=%v", p.err, tracker.Activo("lead-1"))
	}
}

func TestS2MessageValidationAndAudioPrivacy(t *testing.T) {
	large := base64.StdEncoding.EncodeToString(make([]byte, 2*1024*1024+1))
	cases := []struct{ name, body, code string }{
		{"text limit", `{"tipo":"TEXTO","texto":"` + strings.Repeat("x", 2001) + `"}`, "VALIDACION"},
		{"audio mime", `{"tipo":"AUDIO","audio_base64":"aA==","mime":"audio/wav","duracion_s":2}`, "AUDIO_INVALIDO"},
		{"audio duration", `{"tipo":"AUDIO","audio_base64":"aA==","mime":"audio/webm","duracion_s":61}`, "AUDIO_INVALIDO"},
		{"audio decoded size", `{"tipo":"AUDIO","audio_base64":"` + large + `","mime":"audio/webm","duracion_s":2}`, "AUDIO_INVALIDO"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := newAsyncProcessor()
			tracker := newTracker(p)
			defer tracker.Cerrar()
			response := postMessage(messageRouter(asyncLeadRepo(), tracker), tc.body)
			if response.Code != 400 && tc.code == "VALIDACION" || response.Code != 422 && tc.code == "AUDIO_INVALIDO" {
				t.Fatalf("status=%d body=%s", response.Code, response.Body)
			}
			if !strings.Contains(response.Body.String(), `"codigo":"`+tc.code+`"`) || tracker.Activo("lead-1") {
				t.Fatalf("body=%s active=%v", response.Body, tracker.Activo("lead-1"))
			}
		})
	}
	p := newAsyncProcessor()
	tracker := newTracker(p)
	defer tracker.Cerrar()
	body := `{"tipo":"AUDIO","audio_base64":"c2Vuc2l0aXZl","mime":"audio/webm","duracion_s":2}`
	response := postMessage(messageRouter(asyncLeadRepo(), tracker), body)
	if response.Code != 202 || strings.Contains(response.Body.String(), "c2Vuc2l0aXZl") {
		t.Fatalf("audio response=%d %s", response.Code, response.Body)
	}
}
