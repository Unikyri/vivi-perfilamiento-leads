package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

type planTestState struct {
	lead      *domain.Lead
	plan      *domain.PlanNutricion
	createErr error
	saveErr   error
	appendErr error
	order     *[]string
	messages  []domain.Mensaje
	creates   int
}
type planLeadDouble struct{ *planTestState }
type planRepoDouble struct{ *planTestState }

func (f *planLeadDouble) Crear(context.Context, *domain.Lead) error { return nil }
func (f *planLeadDouble) PorID(context.Context, string) (*domain.Lead, error) {
	if f.lead == nil {
		return nil, &NotFoundError{Resource: "lead", ID: "missing"}
	}
	copy := *f.lead
	return &copy, nil
}
func (f *planLeadDouble) Guardar(_ context.Context, lead *domain.Lead) error {
	*f.order = append(*f.order, "save")
	if f.saveErr != nil {
		err := f.saveErr
		f.saveErr = nil
		return err
	}
	f.lead = lead
	return nil
}
func (f *planLeadDouble) Listar(context.Context, FiltroLeads) ([]*domain.Lead, error) {
	return nil, nil
}
func (f *planLeadDouble) AgregarMensaje(_ context.Context, message *domain.Mensaje) error {
	*f.order = append(*f.order, "message")
	if f.appendErr != nil {
		return f.appendErr
	}
	f.messages = append(f.messages, *message)
	return nil
}
func (f *planLeadDouble) Conversacion(context.Context, string) ([]domain.Mensaje, error) {
	return f.messages, nil
}

func (f *planRepoDouble) Crear(_ context.Context, plan *domain.PlanNutricion) error {
	*f.order = append(*f.order, "create")
	f.creates++
	if f.createErr != nil {
		return f.createErr
	}
	copy := *plan
	copy.Hitos = append([]domain.Hito(nil), plan.Hitos...)
	f.plan = &copy
	return nil
}
func (f *planRepoDouble) PorLead(context.Context, string) (*domain.PlanNutricion, error) {
	if f.plan == nil {
		return nil, &NotFoundError{Resource: "plan", ID: "missing"}
	}
	copy := *f.plan
	copy.Hitos = append([]domain.Hito(nil), f.plan.Hitos...)
	return &copy, nil
}
func (f *planRepoDouble) Guardar(_ context.Context, plan *domain.PlanNutricion) error {
	*f.order = append(*f.order, "plan-save")
	f.plan = plan
	return nil
}
func (f *planRepoDouble) HitosVencidos(context.Context, time.Time) ([]HitoConPlan, error) {
	return nil, nil
}
func (f *planRepoDouble) MarcarHito(context.Context, string, domain.EstadoHito) error { return nil }

func newPlanLeadDouble() (*planTestState, *GestionarPlan) {
	order := []string{}
	state := &planTestState{lead: &domain.Lead{LeadID: "lead-1", Ruta: domain.RutaNutricion, Estado: domain.EstadoLeadCalificado, Capacidad: &domain.Capacidad{PresupuestoMax: 100}}, order: &order}
	uc := &GestionarPlan{Leads: &planLeadDouble{state}, Planes: &planRepoDouble{state}, Reloj: NuevoRelojFake(time.Date(2026, 1, 2, 3, 0, 0, 0, time.UTC)), IDs: NuevoIDFake("id"), Calendario: []domain.EventoCalendario{{Tipo: "PRIMA", Fecha: "--06-30"}}}
	return state, uc
}
func TestGestionarPlanCrearValidationHasNoSideEffects(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(EntradaCrearPlan, *domain.Lead) EntradaCrearPlan
	}{
		{"missing capacity", func(in EntradaCrearPlan, lead *domain.Lead) EntradaCrearPlan { lead.Capacidad = nil; return in }},
		{"non-positive target", func(in EntradaCrearPlan, _ *domain.Lead) EntradaCrearPlan { in.PrecioObjetivo = 0; return in }},
		{"invalid frequency", func(in EntradaCrearPlan, _ *domain.Lead) EntradaCrearPlan { in.Frecuencia = "WEEKLY"; return in }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state, uc := newPlanLeadDouble()
			input := tc.edit(EntradaCrearPlan{LeadID: "lead-1", Consintio: true, Frecuencia: "MENSUAL", PrecioObjetivo: 200}, state.lead)
			_, err := uc.CrearPlan(context.Background(), input)
			if !errors.Is(err, ErrValidacion) || state.creates != 0 || len(state.messages) != 0 || state.lead.Estado != domain.EstadoLeadCalificado {
				t.Fatalf("err=%v creates=%d messages=%d state=%s", err, state.creates, len(state.messages), state.lead.Estado)
			}
		})
	}
}
func TestGestionarPlanCrearConsentAndTimestamp(t *testing.T) {
	state, uc := newPlanLeadDouble()
	plan, err := uc.CrearPlan(context.Background(), EntradaCrearPlan{LeadID: "lead-1", Consintio: true, Frecuencia: "MENSUAL", PrecioObjetivo: 200})
	if err != nil || plan == nil || plan.Estado != domain.EstadoPlanActivo || plan.MetaMonto != 100 || plan.ConsentimientoEn == nil || !plan.ConsentimientoEn.Equal(uc.Reloj.Ahora()) || state.lead.Estado != domain.EstadoLeadEnNutricion {
		t.Fatalf("plan=%#v err=%v lead=%s", plan, err, state.lead.Estado)
	}
	if !reflect.DeepEqual(*state.order, []string{"create", "save"}) {
		t.Fatalf("order=%v", *state.order)
	}
	state, uc = newPlanLeadDouble()
	plan, err = uc.CrearPlan(context.Background(), EntradaCrearPlan{LeadID: "lead-1", Consintio: true, Frecuencia: "MENSUAL", PrecioObjetivo: 50})
	if err != nil || plan == nil || plan.MetaMonto != 0 || state.creates != 1 || len(state.messages) != 0 || state.lead.Estado != domain.EstadoLeadEnNutricion {
		t.Fatalf("below budget plan=%#v err=%v creates=%d state=%s", plan, err, state.creates, state.lead.Estado)
	}
}

func TestGestionarPlanCrearNoConsentReminderAndFailure(t *testing.T) {
	for _, wantErr := range []error{nil, errors.New("append failed")} {
		state, uc := newPlanLeadDouble()
		state.appendErr = wantErr
		plan, err := uc.CrearPlan(context.Background(), EntradaCrearPlan{LeadID: "lead-1", Consintio: false, Frecuencia: "MENSUAL", PrecioObjetivo: 200})
		if (wantErr == nil && (err != nil || plan != nil)) || (wantErr != nil && err == nil) || state.creates != 0 || len(*state.order) != 1 || (wantErr == nil && len(state.messages) != 1) {
			t.Fatalf("wantErr=%v plan=%#v err=%v order=%v messages=%d", wantErr, plan, err, *state.order, len(state.messages))
		}
	}
}

func TestGestionarPlanCrearFailureOrderAndRetryReuse(t *testing.T) {
	state, uc := newPlanLeadDouble()
	state.createErr = errors.New("create failed")
	_, err := uc.CrearPlan(context.Background(), EntradaCrearPlan{LeadID: "lead-1", Consintio: true, Frecuencia: "MENSUAL", PrecioObjetivo: 200})
	if err == nil || state.lead.Estado != domain.EstadoLeadCalificado || !reflect.DeepEqual(*state.order, []string{"create"}) {
		t.Fatalf("create failure err=%v state=%s order=%v", err, state.lead.Estado, *state.order)
	}
	state.createErr = nil
	state.saveErr = errors.New("save failed")
	_, err = uc.CrearPlan(context.Background(), EntradaCrearPlan{LeadID: "lead-1", Consintio: true, Frecuencia: "MENSUAL", PrecioObjetivo: 200})
	if err == nil || state.creates != 2 {
		t.Fatalf("save failure err=%v creates=%d", err, state.creates)
	}
	_, err = uc.CrearPlan(context.Background(), EntradaCrearPlan{LeadID: "lead-1", Consintio: true, Frecuencia: "MENSUAL", PrecioObjetivo: 200})
	if err != nil || state.creates != 2 || state.lead.Estado != domain.EstadoLeadEnNutricion || !reflect.DeepEqual(*state.order, []string{"create", "create", "save", "save"}) {
		t.Fatalf("retry err=%v creates=%d state=%s order=%v", err, state.creates, state.lead.Estado, *state.order)
	}
}

func TestGestionarPlanPausarOrderingAndIdempotency(t *testing.T) {
	state, uc := newPlanLeadDouble()
	state.lead.Estado = domain.EstadoLeadEnNutricion
	state.plan = &domain.PlanNutricion{PlanID: "plan-1", LeadID: "lead-1", Estado: domain.EstadoPlanActivo}
	if err := uc.PausarPlan(context.Background(), "lead-1"); err != nil {
		t.Fatal(err)
	}
	if state.plan.Estado != domain.EstadoPlanPausado || state.lead.Estado != domain.EstadoLeadPausado || len(state.messages) != 1 {
		t.Fatalf("plan=%s lead=%s messages=%d", state.plan.Estado, state.lead.Estado, len(state.messages))
	}
	if !reflect.DeepEqual(*state.order, []string{"plan-save", "save", "message"}) {
		t.Fatalf("order=%v", *state.order)
	}
	before := append([]string(nil), *state.order...)
	if err := uc.PausarPlan(context.Background(), "lead-1"); err != nil || !reflect.DeepEqual(*state.order, before) {
		t.Fatalf("repeat err=%v order=%v", err, *state.order)
	}
	state, uc = newPlanLeadDouble()
	state.lead.Estado = domain.EstadoLeadPausado
	state.plan = &domain.PlanNutricion{PlanID: "plan-1", LeadID: "lead-1", Estado: domain.EstadoPlanActivo}
	if err := uc.PausarPlan(context.Background(), "lead-1"); err != nil || state.plan.Estado != domain.EstadoPlanPausado || len(state.messages) != 0 || !reflect.DeepEqual(*state.order, []string{"plan-save"}) {
		t.Fatalf("paused lead err=%v plan=%s messages=%d order=%v", err, state.plan.Estado, len(state.messages), *state.order)
	}
}

func TestGestionarPlanPausarMissingAndFailuresHaveNoFarewell(t *testing.T) {
	for _, tc := range []struct {
		name      string
		lead      domain.EstadoLead
		plan      bool
		saveErr   error
		wantErr   bool
		wantOrder []string
	}{
		{"missing plan", domain.EstadoLeadEnNutricion, false, nil, false, nil},
		{"illegal transition", domain.EstadoLeadCalificado, true, nil, true, []string{"plan-save"}},
		{"lead save failure", domain.EstadoLeadEnNutricion, true, errors.New("save failed"), true, []string{"plan-save", "save"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state, uc := newPlanLeadDouble()
			state.lead.Estado = tc.lead
			if tc.plan {
				state.plan = &domain.PlanNutricion{PlanID: "plan-1", LeadID: "lead-1", Estado: domain.EstadoPlanActivo}
			}
			state.saveErr = tc.saveErr
			err := uc.PausarPlan(context.Background(), "lead-1")
			if (err != nil) != tc.wantErr || len(state.messages) != 0 || len(*state.order) != len(tc.wantOrder) || (len(tc.wantOrder) > 0 && !reflect.DeepEqual(*state.order, tc.wantOrder)) {
				t.Fatalf("err=%v messages=%d order=%v", err, len(state.messages), *state.order)
			}
		})
	}
}
