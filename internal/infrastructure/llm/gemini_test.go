package llm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

func adapterJSON() string {
	return `{"campos_extraidos":[],"intencion":{"nivel":"ALTA","confianza":"MEDIA","senales":["compra"]},"respuesta":"Listo","accion":"CONTINUAR"}`
}
func wireResponse(text string) string {
	b, _ := json.Marshal(text)
	return fmt.Sprintf(`{"candidates":[{"content":{"parts":[{"text":%s}]}}]}`, b)
}
func adapterInput() usecase.EntradaTurno {
	return usecase.EntradaTurno{MensajeUsuario: "hola", NumerosDelMotor: map[string]int64{}}
}

func TestGeminiRequestAndResponse(t *testing.T) {
	var calls int
	client := fakeDoer(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != http.MethodPost || r.URL.Path != "/v1beta/models/test-model:generateContent" || r.Header.Get("X-Goog-Api-Key") != "secret-key" || r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("wire: %s %s", r.URL, r.Header)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["generationConfig"].(map[string]any)["responseMimeType"] != "application/json" || body["systemInstruction"] == nil || body["contents"] == nil {
			t.Fatal("missing Gemini structured fields")
		}
		return response(200, wireResponse(adapterJSON())), nil
	})
	out, err := NewGeminiProvider("secret-key", WithGeminiHTTPDoer(client), WithGeminiBaseURL("https://example.test"), WithGeminiModel("test-model")).GenerarTurno(context.Background(), adapterInput())
	if err != nil || out.Accion != usecase.AccionContinuar || calls != 1 {
		t.Fatalf("out=%#v err=%v calls=%d", out, err, calls)
	}
}

func TestGeminiMalformedRetry(t *testing.T) {
	calls := 0
	client := fakeDoer(func(*http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return response(200, wireResponse("{")), nil
		}
		return response(200, wireResponse(adapterJSON())), nil
	})
	_, err := NewGeminiProvider("key", WithGeminiHTTPDoer(client), WithGeminiBaseURL("http://test")).GenerarTurno(context.Background(), adapterInput())
	if err != nil || calls != 2 {
		t.Fatalf("err=%v calls=%d", err, calls)
	}
}

func TestGeminiAudioAndStatusRedaction(t *testing.T) {
	raw := []byte("audio-bytes")
	client := fakeDoer(func(r *http.Request) (*http.Response, error) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		parts := body["contents"].([]any)[0].(map[string]any)["parts"].([]any)
		inline := parts[1].(map[string]any)["inlineData"].(map[string]any)
		if inline["mimeType"] != "audio/ogg" || inline["data"] != base64.StdEncoding.EncodeToString(raw) {
			t.Fatal("audio wire shape")
		}
		return response(429, "provider-body-with-secret-key"), nil
	})
	_, err := NewGeminiProvider("secret-key", WithGeminiHTTPDoer(client), WithGeminiBaseURL("http://test")).ProcesarAudio(context.Background(), usecase.Audio{Base64: base64.StdEncoding.EncodeToString(raw), MIME: "audio/ogg"}, adapterInput())
	if ErrorKindOf(err) != KindRateLimit || strings.Contains(err.Error(), "secret") || strings.Contains(err.Error(), "provider-body") {
		t.Fatalf("err=%v", err)
	}
}
