package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

type coreLLM struct {
	out                   SalidaTurno
	err                   error
	textCalls, audioCalls int
	last                  EntradaTurno
	lastAudio             Audio
}

func (f *coreLLM) GenerarTurno(_ context.Context, in EntradaTurno) (SalidaTurno, error) {
	f.textCalls++
	f.last = in
	return f.out, f.err
}
func (f *coreLLM) ProcesarAudio(_ context.Context, audio Audio, in EntradaTurno) (SalidaTurno, error) {
	f.audioCalls++
	f.lastAudio, f.last = audio, in
	return f.out, f.err
}
func (f *coreLLM) Nombre() string { return "core-test" }

func coreLead() *domain.Lead {
	return &domain.Lead{LeadID: "lead-1", Estado: domain.EstadoLeadPerfilando, Perfil: domain.Perfil{
		"ingreso_hogar": {Valor: int64(3000000), Fuente: domain.FuenteCampoVerificadoBase},
	}}
}
func coreUC(llm LLMProvider) (*ProcesarMensaje, *LeadRepoFake) {
	repo := NuevoLeadRepoFake()
	_ = repo.Crear(context.Background(), coreLead())
	return &ProcesarMensaje{Leads: repo, LLM: llm, IDs: NuevoIDFake("msg"), Reloj: NuevoRelojFake(time.Unix(100, 0))}, repo
}

func TestSiguientePreguntaCore(t *testing.T) {
	perfil := domain.Perfil{"ingreso_hogar": {Valor: int64(1), Fuente: domain.FuenteCampoVerificadoBase}}
	if got := SiguienteMejorPregunta(perfil); got != "recursos_propios" {
		t.Fatalf("next=%q", got)
	}
	perfil["recursos_propios"] = domain.CampoPerfil{Valor: int64(2)}
	perfil["zona_deseada"] = domain.CampoPerfil{Valor: "Norte"}
	if !PerfilEstaCompleto(perfil) || SiguienteMejorPregunta(perfil) != "plazo_compra_meses" {
		t.Fatalf("completion/question failed")
	}
}

func TestProcesarMensajeCore(t *testing.T) {
	cases := []struct {
		name                string
		input               EntradaMensaje
		providerErr         error
		wantText, wantAudio int
		wantErr             error
	}{
		{"blank text", EntradaMensaje{LeadID: "lead-1", Tipo: domain.TipoMensajeTexto}, nil, 0, 0, ErrValidacion},
		{"invalid audio", EntradaMensaje{LeadID: "lead-1", Tipo: domain.TipoMensajeAudio, Audio: &Audio{Base64: "?", MIME: "audio/webm", DuracionS: 2}}, nil, 0, 0, ErrAudioInvalido},
		{"provider error", EntradaMensaje{LeadID: "lead-1", Tipo: domain.TipoMensajeTexto, Texto: "hola"}, errors.New("provider"), 1, 0, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			llm := &coreLLM{err: tc.providerErr}
			uc, repo := coreUC(llm)
			err := uc.Ejecutar(context.Background(), tc.input)
			if tc.providerErr != nil {
				if !errors.Is(err, tc.providerErr) {
					t.Fatalf("err=%v", err)
				}
			} else if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("err=%v", err)
			}
			if llm.textCalls != tc.wantText || llm.audioCalls != tc.wantAudio {
				t.Fatalf("calls text/audio=%d/%d", llm.textCalls, llm.audioCalls)
			}
			conversation, _ := repo.Conversacion(context.Background(), "lead-1")
			if tc.providerErr != nil || tc.wantErr != nil {
				if len(conversation) != 0 {
					t.Fatalf("unexpected writes=%d", len(conversation))
				}
			}
		})
	}
}

func TestProcesarMensajeCoreAudioAndVerifiedFields(t *testing.T) {
	llm := &coreLLM{out: SalidaTurno{Respuesta: "¿En qué zona?", CamposExtraidos: []CampoExtraido{
		{Campo: "ingreso_hogar", Valor: int64(9), Fuente: domain.FuenteCampoVerificadoBase, Confianza: 1},
		{Campo: "zona_deseada", Valor: "Norte", Fuente: domain.FuenteCampoDeclarado, Confianza: .9},
		{Campo: "inventado", Valor: "ignore", Fuente: domain.FuenteCampoDeclarado, Confianza: .9},
	}}}
	uc, repo := coreUC(llm)
	input := EntradaMensaje{LeadID: "lead-1", Tipo: domain.TipoMensajeAudio, Texto: "quiero comprar", Audio: &Audio{Base64: "aGVsbG8=", MIME: "audio/webm", DuracionS: 2}}
	if err := uc.Ejecutar(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	if llm.audioCalls != 1 || llm.textCalls != 0 || llm.last.MensajeUsuario != input.Texto {
		t.Fatalf("dispatch=%d/%d turn=%q", llm.textCalls, llm.audioCalls, llm.last.MensajeUsuario)
	}
	if llm.lastAudio.Base64 == "" {
		t.Fatal("audio was not dispatched")
	}
	lead, _ := repo.PorID(context.Background(), "lead-1")
	if lead.Perfil["ingreso_hogar"].Valor != int64(3000000) || lead.Perfil["zona_deseada"].Valor != "Norte" {
		t.Fatalf("profile=%v", lead.Perfil)
	}
	if _, ok := lead.Perfil["inventado"]; ok {
		t.Fatal("unknown field merged")
	}
	messages, _ := repo.Conversacion(context.Background(), "lead-1")
	if len(messages) != 2 || messages[0].Autor != domain.AutorMensajeLead || messages[0].Texto != input.Texto || messages[0].Adjunto["audio_original"] != true || messages[1].Texto != llm.out.Respuesta {
		t.Fatalf("messages=%+v", messages)
	}
}
