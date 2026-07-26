package usecase

import (
	"context"
	"errors"
	"reflect"
	"strings"
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
	reloj := NuevoRelojFake(time.Unix(100, 0))
	saludo := &SaludarLead{Leads: repo, IDs: NuevoIDFake("farewell"), Reloj: reloj}
	return &ProcesarMensaje{Leads: repo, LLM: llm, IDs: NuevoIDFake("msg"), Reloj: reloj, Saludo: saludo}, repo
}

func TestSiguientePreguntaCore(t *testing.T) {
	perfil := domain.Perfil{"ingreso_hogar": {Fuente: domain.FuenteCampoVerificadoBase}}
	if got := SiguienteMejorPregunta(perfil); got != "recursos_propios" {
		t.Fatalf("verified nil field must be skipped, next=%q", got)
	}
	perfil["ingreso_hogar"] = domain.CampoPerfil{Valor: int64(1), Fuente: domain.FuenteCampoVerificadoBase}
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

func TestProcesarMensajeAudioProviderErrorHasNoWrites(t *testing.T) {
	providerErr := errors.New("audio provider unavailable")
	llm := &coreLLM{err: providerErr}
	uc, repo := coreUC(llm)
	input := EntradaMensaje{
		LeadID: "lead-1", Tipo: domain.TipoMensajeAudio, Texto: "quiero comprar",
		Audio: &Audio{Base64: "aGVsbG8=", MIME: "audio/webm", DuracionS: 2},
	}
	if err := uc.Ejecutar(context.Background(), input); !errors.Is(err, providerErr) {
		t.Fatalf("err=%v", err)
	}
	if llm.audioCalls != 1 || llm.textCalls != 0 {
		t.Fatalf("audio/text calls=%d/%d", llm.audioCalls, llm.textCalls)
	}
	conversation, _ := repo.Conversacion(context.Background(), "lead-1")
	if len(conversation) != 0 {
		t.Fatalf("provider failure wrote %d messages", len(conversation))
	}
}

type edgeRepo struct {
	*LeadRepoFake
	failAdd, addCalls int
	failSave          bool
}

func (r *edgeRepo) AgregarMensaje(ctx context.Context, m *domain.Mensaje) error {
	r.addCalls++
	if r.failAdd == r.addCalls {
		return errors.New("add failure")
	}
	return r.LeadRepoFake.AgregarMensaje(ctx, m)
}
func (r *edgeRepo) Guardar(ctx context.Context, l *domain.Lead) error {
	if r.failSave {
		return errors.New("cas failure")
	}
	return r.LeadRepoFake.Guardar(ctx, l)
}

type edgeBus struct{ events []Evento }

func (b *edgeBus) Publicar(_ context.Context, e Evento)          { b.events = append(b.events, e) }
func (*edgeBus) Suscribir(string, func(context.Context, Evento)) {}
func edgeSetup(out SalidaTurno) (*ProcesarMensaje, *edgeRepo, *coreLLM, *edgeBus) {
	r := &edgeRepo{LeadRepoFake: NuevoLeadRepoFake()}
	_ = r.Crear(context.Background(), coreLead())
	b, l := &edgeBus{}, &coreLLM{out: out}
	reloj := NuevoRelojFake(time.Unix(100, 0))
	saludo := &SaludarLead{Leads: r, IDs: NuevoIDFake("farewell"), Reloj: reloj}
	return &ProcesarMensaje{Leads: r, LLM: l, IDs: NuevoIDFake("edge"), Bus: b, Reloj: reloj, Saludo: saludo}, r, l, b
}
func edgeInput() EntradaMensaje {
	return EntradaMensaje{LeadID: "lead-1", Tipo: domain.TipoMensajeTexto, Texto: "respuesta"}
}

func TestProcesarMensajeEdges(t *testing.T) {
	t.Run("two pass normalization and protected motor", func(t *testing.T) {
		uc, r, llm, _ := edgeSetup(SalidaTurno{Respuesta: "ok", CamposExtraidos: []CampoExtraido{
			{Campo: "ingreso_hogar", Valor: int64(9), Fuente: domain.FuenteCampoDeclarado, Confianza: 1},
			{Campo: "recursos_propios", Valor: float64(200), Fuente: domain.FuenteCampoDeclarado, Confianza: .8},
			{Campo: "zona_deseada", Valor: "Norte", Fuente: domain.FuenteCampoVerificadoBase, Confianza: .8},
			{Campo: "tipo_hogar", Valor: "APTO", Fuente: domain.FuenteCampoVerificadoBase, Confianza: .8},
			{Campo: "reporte_credito", Valor: true, Fuente: domain.FuenteCampoDeclarado, Confianza: 1},
			{Campo: "desconocido", Valor: true, Fuente: domain.FuenteCampoDeclarado, Confianza: 1},
		}})
		lead, _ := r.PorID(context.Background(), "lead-1")
		lead.Capacidad = &domain.Capacidad{PresupuestoMax: 4100000, CreditoMax: 2800000, SubsidioAplicable: 900000, RecursosPropios: 400000}
		if err := r.Guardar(context.Background(), lead); err != nil {
			t.Fatal(err)
		}
		if err := uc.Ejecutar(context.Background(), edgeInput()); err != nil {
			t.Fatal(err)
		}
		wantMotor := map[string]int64{
			"presupuesto_max":    4100000,
			"credito_max":        2800000,
			"subsidio_aplicable": 900000,
			"recursos_propios":   400000,
		}
		for campo, want := range wantMotor {
			if got := llm.last.NumerosDelMotor[campo]; got != want {
				t.Fatalf("motor[%q]=%d, want %d; captured=%v", campo, got, want, llm.last.NumerosDelMotor)
			}
		}
		l, _ := r.PorID(context.Background(), "lead-1")
		if l.Perfil["ingreso_hogar"].Valor != int64(3000000) || l.Perfil["recursos_propios"].Valor != int64(200) || l.Perfil["zona_deseada"].Fuente != domain.FuenteCampoDeclarado || l.Perfil["reporte_credito"].Valor != true || l.Capacidad.RecursosPropios != 200 {
			t.Fatalf("profile/capacity=%+v/%+v", l.Perfil, l.Capacidad)
		}
	})
	t.Run("invalid recognized output is side effect free", func(t *testing.T) {
		for _, fields := range [][]CampoExtraido{{{Campo: "edad", Valor: "x", Fuente: domain.FuenteCampoDeclarado}}, {{Campo: "edad", Valor: 2, Fuente: "bad"}}, {{Campo: "edad", Valor: 1, Fuente: domain.FuenteCampoDeclarado}, {Campo: "edad", Valor: 2, Fuente: domain.FuenteCampoDeclarado}}, {{Campo: "reporte_credito", Valor: "true", Fuente: domain.FuenteCampoDeclarado}}} {
			uc, r, _, _ := edgeSetup(SalidaTurno{CamposExtraidos: fields})
			before, _ := r.PorID(context.Background(), "lead-1")
			if err := uc.Ejecutar(context.Background(), edgeInput()); !errors.Is(err, ErrSalidaTurnoInvalida) {
				t.Fatalf("err=%v", err)
			}
			after, _ := r.PorID(context.Background(), "lead-1")
			messages, _ := r.Conversacion(context.Background(), "lead-1")
			if len(messages) != 0 || len(after.Perfil) != len(before.Perfil) || after.Estado != before.Estado || after.Capacidad != before.Capacidad {
				t.Fatalf("mutation/messages: before=%+v after=%+v writes=%d", before, after, len(messages))
			}
		}
	})
	t.Run("unintelligible audio stores response only", func(t *testing.T) {
		uc, r, l, b := edgeSetup(SalidaTurno{Respuesta: "repite", Accion: AccionAudioInint, CamposExtraidos: []CampoExtraido{{Campo: "zona_deseada", Valor: "Norte", Fuente: domain.FuenteCampoDeclarado}}})
		before, _ := r.PorID(context.Background(), "lead-1")
		before.Capacidad = &domain.Capacidad{PresupuestoMax: 7}
		_ = r.Guardar(context.Background(), before)
		in := edgeInput()
		in.Tipo, in.Audio = domain.TipoMensajeAudio, &Audio{Base64: "aA==", MIME: "audio/webm", DuracionS: 2}
		if err := uc.Ejecutar(context.Background(), in); err != nil {
			t.Fatal(err)
		}
		lead, _ := r.PorID(context.Background(), "lead-1")
		messages, _ := r.Conversacion(context.Background(), "lead-1")
		if l.audioCalls != 1 || len(messages) != 1 || messages[0].Autor != domain.AutorMensajeVivi || lead.Capacidad.PresupuestoMax != 7 || lead.Perfil["zona_deseada"].Valor != nil || len(b.events) != 0 {
			t.Fatalf("audio state/messages/events=%+v/%+v/%d", lead, messages, len(b.events))
		}
	})
	t.Run("affiliate cap completes and publishes once", func(t *testing.T) {
		uc, r, l, b := edgeSetup(SalidaTurno{Respuesta: "listo", Accion: AccionContinuar})
		lead, _ := r.PorID(context.Background(), "lead-1")
		lead.Afiliado = true
		_ = r.Guardar(context.Background(), lead)
		for i := 0; i < 3; i++ {
			_ = r.AgregarMensaje(context.Background(), &domain.Mensaje{LeadID: "lead-1", Autor: domain.AutorMensajeLead, CreadoEn: time.Unix(int64(i), 0)})
		}
		if err := uc.Ejecutar(context.Background(), edgeInput()); err != nil {
			t.Fatal(err)
		}
		got, _ := r.PorID(context.Background(), "lead-1")
		if got.Estado != domain.EstadoLeadCalificado || len(b.events) != 1 || l.textCalls != 1 {
			t.Fatalf("state/events/calls=%s/%d/%d", got.Estado, len(b.events), l.textCalls)
		}
		if err := uc.Ejecutar(context.Background(), edgeInput()); !errors.Is(err, ErrLeadNoPerfilando) {
			t.Fatalf("post completion=%v", err)
		}
	})
	t.Run("non affiliate cap and bounded history", func(t *testing.T) {
		uc, r, l, _ := edgeSetup(SalidaTurno{Respuesta: "listo", Accion: AccionContinuar})
		for i := 0; i < 8; i++ {
			a := domain.AutorMensajeVivi
			if i < 5 {
				a = domain.AutorMensajeLead
			}
			_ = r.AgregarMensaje(context.Background(), &domain.Mensaje{LeadID: "lead-1", Autor: a, Texto: string(rune('a' + i)), CreadoEn: time.Unix(int64(i), 0)})
		}
		if err := uc.Ejecutar(context.Background(), edgeInput()); err != nil || l.textCalls != 1 || len(l.last.HistorialReciente) != 7 {
			t.Fatalf("err/calls/history=%v/%d/%d", err, l.textCalls, len(l.last.HistorialReciente))
		}
		blocked, rr, ll, _ := edgeSetup(SalidaTurno{Accion: AccionContinuar})
		for i := 0; i < 6; i++ {
			_ = rr.AgregarMensaje(context.Background(), &domain.Mensaje{LeadID: "lead-1", Autor: domain.AutorMensajeLead})
		}
		if err := blocked.Ejecutar(context.Background(), edgeInput()); !errors.Is(err, ErrLimiteTurnos) || ll.textCalls != 0 {
			t.Fatalf("post cap err/calls=%v/%d", err, ll.textCalls)
		}
	})
	t.Run("failure prefixes publish no event", func(t *testing.T) {
		for _, fail := range []int{1, 2} {
			uc, r, _, b := edgeSetup(SalidaTurno{Respuesta: "ok", Accion: AccionPerfilCompleto})
			r.failAdd = fail
			err := uc.Ejecutar(context.Background(), edgeInput())
			if err == nil || !strings.Contains(err.Error(), "guardar") || len(b.events) != 0 {
				t.Fatalf("fail=%d err/events=%v/%d", fail, err, len(b.events))
			}
		}
		uc, r, _, b := edgeSetup(SalidaTurno{Respuesta: "ok", Accion: AccionPerfilCompleto})
		r.failSave = true
		err := uc.Ejecutar(context.Background(), edgeInput())
		if err == nil || !strings.Contains(err.Error(), "guardar lead") || len(b.events) != 0 {
			t.Fatalf("cas err/events=%v/%d", err, len(b.events))
		}
	})
}

func TestProcesarMensajeUsesAcceptedMetadata(t *testing.T) {
	uc, repo := coreUC(&coreLLM{out: SalidaTurno{Respuesta: "ok"}})
	received := time.Date(2026, 7, 26, 7, 0, 0, 0, time.UTC)
	input := EntradaMensaje{LeadID: "lead-1", Tipo: domain.TipoMensajeTexto, Texto: "hola", MensajeID: "accepted-1", RecibidoEn: received}
	if err := uc.Ejecutar(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	messages, _ := repo.Conversacion(context.Background(), "lead-1")
	found := false
	for _, message := range messages {
		if message.MensajeID == "accepted-1" && message.CreadoEn.Equal(received) {
			found = true
		}
	}
	if !found {
		t.Fatalf("message=%+v", messages)
	}
}

func TestProcesarMensajeConsentDenialIsTerminalAndProfileSafe(t *testing.T) {
	llm := &coreLLM{out: SalidaTurno{
		Accion: AccionConsentimientoNo, Respuesta: "no gracias",
		CamposExtraidos: []CampoExtraido{{Campo: "ingreso_hogar", Valor: int64(99999999), Fuente: domain.FuenteCampoDeclarado, Confianza: 1}},
		Intencion:       domain.Intencion{Nivel: domain.NivelAlta, Confianza: domain.NivelAlta},
	}}
	uc, repo, provider, bus := edgeSetup(llm.out)
	before, _ := repo.PorID(context.Background(), "lead-1")
	before.Perfil["recursos_propios"] = domain.CampoPerfil{Valor: int64(2000000), Fuente: domain.FuenteCampoDeclarado}
	before.Perfil["zona_deseada"] = domain.CampoPerfil{Valor: "Norte", Fuente: domain.FuenteCampoDeclarado}
	if !PerfilEstaCompleto(before.Perfil) {
		t.Fatal("denial fixture must be structurally complete before processing")
	}
	before.Capacidad = &domain.Capacidad{PresupuestoMax: 10, SubsidioAplicable: 4}
	before.Intencion = &domain.Intencion{Nivel: domain.NivelBaja, Confianza: domain.NivelMedia, Senales: []string{"original"}}
	if err := repo.Guardar(context.Background(), before); err != nil {
		t.Fatal(err)
	}
	input := EntradaMensaje{LeadID: "lead-1", Tipo: domain.TipoMensajeTexto, Texto: "No autorizo", MensajeID: "denial-1"}
	if err := uc.Ejecutar(context.Background(), input); err != nil {
		t.Fatal(err)
	}
	after, _ := repo.PorID(context.Background(), "lead-1")
	messages, _ := repo.Conversacion(context.Background(), "lead-1")
	if after.Estado != domain.EstadoLeadDespedido || after.Ruta != domain.RutaDespedida {
		t.Fatalf("terminal route/state=%s/%s", after.Estado, after.Ruta)
	}
	if !reflect.DeepEqual(after.Perfil, before.Perfil) || !reflect.DeepEqual(after.Capacidad, before.Capacidad) || !reflect.DeepEqual(after.Intencion, before.Intencion) {
		t.Fatalf("profile data changed: before=%+v after=%+v", before, after)
	}
	if len(messages) != 2 || messages[0].Texto != input.Texto || messages[0].Autor != domain.AutorMensajeLead || messages[1].Autor != domain.AutorMensajeVivi || !strings.Contains(messages[1].Texto, "Respeto tu decisión") {
		t.Fatalf("denial conversation=%+v", messages)
	}
	if len(bus.events) != 0 || provider.textCalls != 1 {
		t.Fatalf("events/calls=%d/%d", len(bus.events), llm.textCalls)
	}
	for _, event := range bus.events {
		if event.Tipo == EvPerfilCompleto {
			t.Fatal("consent denial published PerfilCompleto")
		}
	}
}

func TestRechazarConsentimientoStopsAfterOrderedWriteFailures(t *testing.T) {
	for _, tc := range []struct {
		name       string
		failAdd    int
		failSave   bool
		wantWrites int
	}{
		{"inbound failure", 1, false, 0},
		{"farewell failure", 2, false, 1},
		{"CAS failure", 0, true, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &edgeRepo{LeadRepoFake: NuevoLeadRepoFake(), failAdd: tc.failAdd, failSave: tc.failSave}
			lead := coreLead()
			if err := r.Crear(context.Background(), lead); err != nil {
				t.Fatal(err)
			}
			clock := NuevoRelojFake(time.Unix(100, 0))
			uc := &SaludarLead{Leads: r, IDs: NuevoIDFake("denial"), Reloj: clock}
			err := uc.RechazarConsentimiento(context.Background(), lead, EntradaMensaje{LeadID: lead.LeadID, Tipo: domain.TipoMensajeTexto, Texto: "No"})
			if err == nil {
				t.Fatal("expected ordered write failure")
			}
			messages, _ := r.Conversacion(context.Background(), lead.LeadID)
			if len(messages) != tc.wantWrites {
				t.Fatalf("messages=%d want %d", len(messages), tc.wantWrites)
			}
		})
	}
}
