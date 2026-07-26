package usecase

import (
	"context"
	"errors"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"reflect"
	"strings"
	"testing"
	"time"
)

type tickState struct {
	lead                        *domain.Lead
	plan                        *domain.PlanNutricion
	due                         []HitoConPlan
	errSend, errAppend, errMark error
	events                      []string
	marked                      []string
	messages                    []domain.Mensaje
	published                   []Evento
	now                         time.Time
}
type tickLeads struct {
	s *tickState
	LeadRepository
}
type tickPlans struct {
	s *tickState
	PlanRepository
}
type tickGateway struct{ s *tickState }
type tickBus struct{ s *tickState }

func (r *tickLeads) PorID(context.Context, string) (*domain.Lead, error) {
	if r.s.lead == nil {
		return nil, &NotFoundError{Resource: "lead", ID: "missing"}
	}
	return r.s.lead, nil
}
func (r *tickLeads) Guardar(_ context.Context, lead *domain.Lead) error {
	r.s.events = append(r.s.events, "lead-save")
	r.s.lead = lead
	return nil
}
func (r *tickLeads) AgregarMensaje(_ context.Context, message *domain.Mensaje) error {
	r.s.events = append(r.s.events, "append")
	if r.s.errAppend != nil {
		return r.s.errAppend
	}
	r.s.messages = append(r.s.messages, *message)
	return nil
}
func (r *tickPlans) PorLead(context.Context, string) (*domain.PlanNutricion, error) {
	return r.s.plan, nil
}
func (r *tickPlans) HitosVencidos(context.Context, time.Time) ([]HitoConPlan, error) {
	return r.s.due, nil
}
func (r *tickPlans) MarcarHito(_ context.Context, id string, state domain.EstadoHito) error {
	r.s.events = append(r.s.events, "mark")
	if r.s.errMark != nil {
		return r.s.errMark
	}
	r.s.marked = append(r.s.marked, id)
	r.s.plan.Hitos[len(r.s.marked)-1].Estado = state
	return nil
}
func (g *tickGateway) Enviar(_ context.Context, _ *domain.Mensaje) error {
	g.s.events = append(g.s.events, "send")
	return g.s.errSend
}
func (b *tickBus) Publicar(_ context.Context, event Evento) {
	b.s.published = append(b.s.published, event)
}
func (b *tickBus) Suscribir(string, func(context.Context, Evento)) {}

type tickClock struct{ now time.Time }

func (c *tickClock) Ahora() time.Time         { return c.now }
func (c *tickClock) FechaSimulada() time.Time { return c.now }
func (c *tickClock) Avanzar(now time.Time)    { c.now = now }
func newTick(hitos ...domain.Hito) (*tickState, *EjecutarHitos) {
	now := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	lead := &domain.Lead{LeadID: "lead-1", Estado: domain.EstadoLeadEnNutricion, Ruta: domain.RutaNutricion}
	plan := &domain.PlanNutricion{PlanID: "plan-1", LeadID: lead.LeadID, Estado: domain.EstadoPlanActivo, MetaMonto: 100, Hitos: append([]domain.Hito(nil), hitos...)}
	due := make([]HitoConPlan, len(hitos))
	for i, h := range hitos {
		due[i] = HitoConPlan{Hito: h, PlanID: plan.PlanID, LeadID: lead.LeadID}
	}
	s := &tickState{lead: lead, plan: plan, due: due, now: now}
	uc := &EjecutarHitos{Leads: &tickLeads{s, nil}, Planes: &tickPlans{s, nil}, Gateway: &tickGateway{s}, Reloj: &tickClock{now: now}, IDs: NuevoIDFake("msg"), Bus: &tickBus{s}}
	return s, uc
}
func TestEjecutarHitosOrdersAndUsesDignifiedPauseCopy(t *testing.T) {
	s, uc := newTick(domain.Hito{HitoID: "h1", Tipo: domain.TipoHitoAhorro, Descripcion: "revisar tu avance", Estado: domain.EstadoHitoPendiente, Monto: int64ptr(25)})
	if _, err := uc.Ejecutar(context.Background(), s.now.Add(-time.Minute)); !errors.Is(err, ErrTiempoSimuladoAtras) {
		t.Fatalf("backward err=%v", err)
	}
	s.lead.Estado = domain.EstadoLeadPausado
	if count, err := uc.Ejecutar(context.Background(), s.now); err != nil || count != 0 {
		t.Fatalf("paused count=%d err=%v", count, err)
	}
	s.lead.Estado = domain.EstadoLeadEnNutricion
	count, err := uc.Ejecutar(context.Background(), s.now.Add(time.Hour))
	if err != nil || count != 1 || !reflect.DeepEqual(s.events, []string{"send", "append", "mark"}) || len(s.marked) != 1 {
		t.Fatalf("count=%d err=%v events=%v marked=%v", count, err, s.events, s.marked)
	}
	text := s.messages[0].Texto
	if !strings.HasSuffix(text, "Si prefieres pausar estos mensajes, responde PAUSAR.") || strings.Contains(strings.ToLower(text), "deuda") || strings.Contains(strings.ToLower(text), "distres") {
		t.Fatalf("text=%q", text)
	}
}
func TestEjecutarHitosContinuesFailuresAndRetriesPending(t *testing.T) {
	s, uc := newTick(
		domain.Hito{HitoID: "h1", Tipo: domain.TipoHitoAhorro, Estado: domain.EstadoHitoPendiente},
		domain.Hito{HitoID: "h2", Tipo: domain.TipoHitoPrima, Estado: domain.EstadoHitoPendiente},
	)
	s.errSend = errors.New("send failed")
	count, err := uc.Ejecutar(context.Background(), s.now)
	if err == nil || count != 0 || len(s.events) != 2 || len(s.marked) != 0 {
		t.Fatalf("count=%d err=%v events=%v marked=%v", count, err, s.events, s.marked)
	}
	s.errSend = nil
	s.errAppend = errors.New("append failed")
	_, err = uc.Ejecutar(context.Background(), s.now)
	if err == nil || len(s.marked) != 0 {
		t.Fatalf("append retry err=%v marked=%v", err, s.marked)
	}
	s.errAppend, s.errMark = nil, errors.New("mark failed")
	_, err = uc.Ejecutar(context.Background(), s.now)
	if err == nil {
		t.Fatalf("mark retry err=%v", err)
	}
}
func TestEjecutarHitosHandoffOnceAndNilBus(t *testing.T) {
	s, uc := newTick(
		domain.Hito{HitoID: "h1", Tipo: domain.TipoHitoAhorro, Estado: domain.EstadoHitoPendiente, Monto: int64ptr(60)},
		domain.Hito{HitoID: "h2", Tipo: domain.TipoHitoAhorro, Estado: domain.EstadoHitoPendiente, Monto: int64ptr(40)},
	)
	count, err := uc.Ejecutar(context.Background(), s.now)
	if err != nil || count != 2 || s.lead.Estado != domain.EstadoLeadPerfilando || len(s.published) != 1 || s.published[0].Tipo != EvPerfilCompleto {
		t.Fatalf("count=%d err=%v state=%s events=%v", count, err, s.lead.Estado, s.published)
	}
	s, uc = newTick(domain.Hito{HitoID: "zero", Tipo: domain.TipoHitoAhorro, Estado: domain.EstadoHitoPendiente, Monto: int64ptr(1)})
	s.plan.MetaMonto = 0
	if count, err := uc.Ejecutar(context.Background(), s.now); err != nil || count != 1 || s.lead.Estado != domain.EstadoLeadEnNutricion || len(s.published) != 0 {
		t.Fatalf("zero meta count=%d err=%v state=%s events=%v", count, err, s.lead.Estado, s.published)
	}
	s, uc = newTick(domain.Hito{HitoID: "re", Tipo: domain.TipoHitoReevaluacion, Estado: domain.EstadoHitoPendiente})
	uc.Bus = nil
	if _, err := uc.Ejecutar(context.Background(), s.now); err != nil || s.lead.Estado != domain.EstadoLeadPerfilando {
		t.Fatalf("nil bus err=%v state=%s", err, s.lead.Estado)
	}
}
func int64ptr(v int64) *int64 { return &v }
