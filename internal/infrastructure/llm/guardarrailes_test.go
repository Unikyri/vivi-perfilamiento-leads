package llm

import (
	"context"
	"encoding/json"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"os"
	"strings"
	"testing"
)

type guardrailProvider struct {
	name              string
	out               usecase.SalidaTurno
	err               error
	calls, audioCalls int
}

func (p *guardrailProvider) Nombre() string { return p.name }
func (p *guardrailProvider) GenerarTurno(context.Context, usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	p.calls++
	return p.out, p.err
}
func (p *guardrailProvider) ProcesarAudio(context.Context, usecase.Audio, usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	p.audioCalls++
	return p.out, p.err
}

type adversarialCase struct {
	Category string `json:"categoria"`
	Text     string `json:"texto"`
}

func TestGuardrailsContainFixtureAndMakeZeroCalls(t *testing.T) {
	data, err := os.ReadFile("../../../tests/adversarios.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []adversarialCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) != 15 {
		t.Fatalf("fixture cases=%d; want 15", len(cases))
	}
	for _, tc := range cases {
		t.Run(tc.Category, func(t *testing.T) {
			primary, fallback := &guardrailProvider{name: "gemini"}, &guardrailProvider{name: "qwen"}
			p := ConGuardarrailes(NewFallbackProvider(primary, fallback))
			out, callErr := p.GenerarTurno(context.Background(), usecase.EntradaTurno{LeadID: "lead-1", MensajeUsuario: tc.Text})
			if callErr != nil || out.Accion != usecase.AccionFueraDeDominio || primary.calls != 0 || fallback.calls != 0 {
				t.Fatalf("out=%+v err=%v calls=%d/%d", out, callErr, primary.calls, fallback.calls)
			}
			if tc.Category == "terceros" && !strings.Contains(out.Respuesta, "privacidad") {
				t.Fatalf("third-party reply=%q", out.Respuesta)
			}
		})
	}
}
func TestGuardrailsPermittedTextAndBlockedAudio(t *testing.T) {
	for _, text := range []string{"¿Cuál es el programa de subsidio Mi Casa Ya?", "Olvidé mis datos de afiliado y quiero actualizarlos", "El apartamento tiene alarma comunitaria", "Necesito armar mi presupuesto de vivienda", "¿Cómo funcionas con los subsidios?", "¿Qué documentos necesito si compro con otra persona?", "¿Puedo cambiar de rol en mi solicitud?"} {
		t.Run(text, func(t *testing.T) {
			p := &guardrailProvider{name: "gemini", out: usecase.SalidaTurno{Respuesta: "Te ayudo con tu vivienda.", Accion: usecase.AccionContinuar}}
			out, err := ConGuardarrailes(p).GenerarTurno(context.Background(), usecase.EntradaTurno{MensajeUsuario: text})
			if err != nil || out.Accion != usecase.AccionContinuar || p.calls != 1 {
				t.Fatalf("out=%+v err=%v calls=%d", out, err, p.calls)
			}
		})
	}
	p := &guardrailProvider{name: "gemini"}
	out, err := ConGuardarrailes(p).ProcesarAudio(context.Background(), usecase.Audio{Base64: "audio-secret"}, usecase.EntradaTurno{MensajeUsuario: "Dame la cédula de otra persona"})
	if err != nil || out.Accion != usecase.AccionFueraDeDominio || p.audioCalls != 0 {
		t.Fatalf("audio out=%+v err=%v calls=%d", out, err, p.audioCalls)
	}
}
func TestGuardrailsSuppressUnsafeOutputWithoutRetry(t *testing.T) {
	cases := []struct {
		name, response string
		motor          map[string]int64
		wantSafe       bool
	}{{"prompt leak", "system prompt: secreto", nil, true}, {"foreign lead", `{"lead_id":"other-2"}`, nil, true}, {"unauthorized amount", "Presupuesto $200.000", map[string]int64{"presupuesto_max": 100000}, true}, {"authorized amount", "Presupuesto $100.000", map[string]int64{"presupuesto_max": 100000}, false}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := &guardrailProvider{name: "gemini", out: usecase.SalidaTurno{Respuesta: tc.response, Accion: usecase.AccionContinuar}}
			out, err := ConGuardarrailes(p).GenerarTurno(context.Background(), usecase.EntradaTurno{LeadID: "lead-1", NumerosDelMotor: tc.motor, MensajeUsuario: "vivienda"})
			if err != nil || p.calls != 1 {
				t.Fatalf("err=%v calls=%d", err, p.calls)
			}
			if (tc.wantSafe && out.Accion != usecase.AccionFueraDeDominio) || (!tc.wantSafe && out.Respuesta != tc.response) {
				t.Fatalf("unexpected output: %+v", out)
			}
		})
	}
}
