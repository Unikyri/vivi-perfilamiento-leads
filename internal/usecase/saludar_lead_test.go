package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

type greetingProvider struct {
	out   SalidaTurno
	err   error
	calls int
	last  EntradaTurno
}

func (p *greetingProvider) GenerarTurno(_ context.Context, in EntradaTurno) (SalidaTurno, error) {
	p.calls++
	p.last = in
	return p.out, p.err
}
func (*greetingProvider) ProcesarAudio(context.Context, Audio, EntradaTurno) (SalidaTurno, error) {
	return SalidaTurno{}, errors.New("audio not expected")
}
func (*greetingProvider) Nombre() string { return "greeting-test" }

func greetingLead(affiliate bool, subsidy int64) *domain.Lead {
	return &domain.Lead{
		LeadID: "lead-1", Nombre: "Ana", Estado: domain.EstadoLeadPerfilando,
		Afiliado: affiliate, Capacidad: &domain.Capacidad{SubsidioAplicable: subsidy},
	}
}

func greetingUC(lead *domain.Lead, provider LLMProvider) (*SaludarLead, *LeadRepoFake) {
	repo := NuevoLeadRepoFake()
	if err := repo.Crear(context.Background(), lead); err != nil {
		panic(err)
	}
	return &SaludarLead{Leads: repo, LLM: provider, IDs: NuevoIDFake("msg"), Reloj: NuevoRelojFake(time.Unix(10, 0))}, repo
}

func TestSaludarLeadUsesOneValidatedProviderDraftOrFallback(t *testing.T) {
	affiliateText := "¡Hola Ana! Tengo una orientación para ti con $52,5M. " + URLPolitica + ". Al continuar autorizas el tratamiento de tus datos. ¿Qué sueñas con comprar?"
	cases := []struct {
		name       string
		lead       *domain.Lead
		provider   LLMProvider
		wantText   string
		wantAmount string
		wantInput  bool
	}{
		{"compliant affiliate draft", greetingLead(true, 52_500_000), &greetingProvider{out: SalidaTurno{Respuesta: affiliateText}}, affiliateText, "$52,5M", true},
		{"provider error uses affiliate fallback", greetingLead(true, 52_500_000), &greetingProvider{err: errors.New("provider unavailable")}, "", "$52,5M", true},
		{"invalid draft uses non affiliate fallback", greetingLead(false, 0), &greetingProvider{out: SalidaTurno{Respuesta: "¿Dos preguntas? ¿Otra?"}}, "", "", true},
		{"nil provider uses deterministic fallback", greetingLead(false, 0), nil, "", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uc, repo := greetingUC(tc.lead, tc.provider)
			if err := uc.Ejecutar(context.Background(), Evento{Tipo: EvLeadNuevo, LeadID: tc.lead.LeadID}); err != nil {
				t.Fatal(err)
			}
			messages, err := repo.Conversacion(context.Background(), tc.lead.LeadID)
			if err != nil || len(messages) != 1 {
				t.Fatalf("messages=%+v err=%v", messages, err)
			}
			if tc.wantText != "" && messages[0].Texto != tc.wantText {
				t.Fatalf("text=%q want %q", messages[0].Texto, tc.wantText)
			}
			if tc.wantText == "" && !strings.Contains(messages[0].Texto, URLPolitica) {
				t.Fatalf("fallback missing policy: %q", messages[0].Texto)
			}
			if !ValidarSaludo(messages[0].Texto, tc.lead.Nombre, tc.lead.Afiliado, tc.wantAmount) {
				t.Fatalf("persisted greeting violates the audience contract: %q", messages[0].Texto)
			}
			if provider, ok := tc.provider.(*greetingProvider); ok && provider.calls != 1 {
				t.Fatalf("provider calls=%d, want 1", provider.calls)
			}
			if provider, ok := tc.provider.(*greetingProvider); ok && provider.last.NumerosDelMotor["subsidio_aplicable"] != tc.lead.Capacidad.SubsidioAplicable {
				t.Fatalf("motor input=%v", provider.last.NumerosDelMotor)
			}
		})
	}
}

func TestFormatoSubsidioAlwaysUsesOneDecimalMillion(t *testing.T) {
	for _, tc := range []struct {
		value int64
		want  string
	}{
		{52_000_000, "$52,0M"},
		{52_500_000, "$52,5M"},
	} {
		t.Run(tc.want, func(t *testing.T) {
			if got := formatoSubsidio(tc.value); got != tc.want {
				t.Fatalf("formatoSubsidio(%d) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
}

func TestValidarSaludoRejectsUnsafeDrafts(t *testing.T) {
	validAffiliate := "¡Hola Ana! " + formatoSubsidio(52_500_000) + " y " + URLPolitica + ". Al continuar autorizas el tratamiento de tus datos. ¿Qué sueñas con comprar?"
	validJob := "¡Hola Ana! " + URLPolitica + ". Al continuar autorizas el tratamiento de tus datos. ¿Cómo está tu situación laboral?"
	cases := []struct {
		name string
		text string
		aff  bool
		amt  string
		want bool
	}{
		{"affiliate accepted", validAffiliate, true, "$52,5M", true},
		{"wrong amount rejected", strings.Replace(validAffiliate, "$52,5M", "$40M", 1), true, "$52,5M", false},
		{"income prompt rejected", strings.Replace(validAffiliate, "¿Qué sueñas", "¿Cuál es tu ingreso y qué sueñas", 1), true, "$52,5M", false},
		{"job copy accepted", validJob, false, "", true},
		{"non affiliate amount rejected", strings.Replace(validJob, URLPolitica, URLPolitica+" por $2M", 1), false, "", false},
		{"two questions rejected", strings.Replace(validJob, "¿Cómo está", "¿Cómo está", 1) + " ¿Otra?", false, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ValidarSaludo(tc.text, "Ana", tc.aff, tc.amt); got != tc.want {
				t.Fatalf("valid=%v want %v text=%q", got, tc.want, tc.text)
			}
		})
	}
}
