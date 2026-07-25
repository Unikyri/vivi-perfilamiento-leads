package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

func qwenWireResponse(text string) string {
	b, _ := json.Marshal(text)
	return `{"choices":[{"message":{"content":` + string(b) + `}}]}`
}

func TestQwenRequestAndResponse(t *testing.T) {
	calls := 0
	client := fakeDoer(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != http.MethodPost || r.URL.Path != "/v1/chat/completions" || r.Header.Get("Authorization") != "Bearer qwen-secret" {
			t.Fatalf("wire: %s %s", r.URL, r.Header)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		messages := body["messages"].([]any)
		if body["model"] != "qwen-test" || body["response_format"].(map[string]any)["type"] != "json_object" || messages[0].(map[string]any)["role"] != "system" || messages[1].(map[string]any)["role"] != "user" {
			t.Fatal("Qwen body shape")
		}
		return response(200, qwenWireResponse(adapterJSON())), nil
	})
	out, err := NewQwenProvider("qwen-secret", "https://example.test/v1", WithQwenHTTPDoer(client), WithQwenModel("qwen-test")).GenerarTurno(context.Background(), adapterInput())
	if err != nil || out.Accion != usecase.AccionContinuar || calls != 1 {
		t.Fatalf("out=%#v err=%v calls=%d", out, err, calls)
	}
}

func TestQwenMalformedRetryAndStatuses(t *testing.T) {
	calls := 0
	client := fakeDoer(func(*http.Request) (*http.Response, error) {
		calls++
		if calls < 3 {
			return response(200, qwenWireResponse("{")), nil
		}
		return response(200, qwenWireResponse(adapterJSON())), nil
	})
	p := NewQwenProvider("key", "http://test", WithQwenHTTPDoer(client))
	_, err := p.GenerarTurno(context.Background(), adapterInput())
	if ErrorKindOf(err) != KindMalformed || calls != 2 {
		t.Fatalf("retry err=%v calls=%d", err, calls)
	}
	for _, status := range []int{429, 500} {
		_, err = NewQwenProvider("qwen-secret", "http://test", WithQwenHTTPDoer(fakeDoer(func(*http.Request) (*http.Response, error) { return response(status, "body-secret"), nil }))).GenerarTurno(context.Background(), adapterInput())
		if ErrorKindOf(err) == "" || strings.Contains(err.Error(), "secret") {
			t.Fatalf("status=%d err=%v", status, err)
		}
	}
}

func TestQwenAudioHasNoTransportOrReroute(t *testing.T) {
	calls := 0
	p := NewQwenProvider("key", "http://qwen", WithQwenHTTPDoer(fakeDoer(func(*http.Request) (*http.Response, error) {
		calls++
		return response(200, qwenWireResponse(adapterJSON())), nil
	})))
	_, err := p.ProcesarAudio(context.Background(), usecase.Audio{Base64: "not-used"}, adapterInput())
	if ErrorKindOf(err) != KindCapability || calls != 0 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}
