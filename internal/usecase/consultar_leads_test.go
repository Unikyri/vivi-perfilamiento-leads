package usecase

import (
	"context"
	"errors"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"testing"
	"time"
)

type planReadFake struct {
	plan *domain.PlanNutricion
	err  error
}

func (f planReadFake) Crear(context.Context, *domain.PlanNutricion) error   { return nil }
func (f planReadFake) Guardar(context.Context, *domain.PlanNutricion) error { return nil }
func (f planReadFake) HitosVencidos(context.Context, time.Time) ([]HitoConPlan, error) {
	return nil, nil
}
func (f planReadFake) MarcarHito(context.Context, string, domain.EstadoHito) error { return nil }
func (f planReadFake) PorLead(context.Context, string) (*domain.PlanNutricion, error) {
	return f.plan, f.err
}

func TestConsultarLeadsFilteredQueueKeepsGlobalCupo(t *testing.T) {
	repo := NuevoLeadRepoFake()
	for _, lead := range []*domain.Lead{
		{LeadID: "advisor", Ruta: domain.RutaAsesor, Prioridad: 1},
		{LeadID: "nutrition", Ruta: domain.RutaNutricion, Prioridad: .5},
	} {
		if err := repo.Crear(context.Background(), lead); err != nil {
			t.Fatal(err)
		}
	}
	route := domain.RutaNutricion
	out, err := (&ConsultarLeads{Leads: repo}).Ejecutar(context.Background(), FiltroLeads{Ruta: &route})
	if err != nil || len(out.Leads) != 1 || out.Leads[0].LeadID != "nutrition" || out.Cupo10.Usados != 1 {
		t.Fatalf("filtered queue=%+v err=%v", out, err)
	}
}

func TestConsultarLeadsDetailIncludesSemaforoAndPlan(t *testing.T) {
	repo := NuevoLeadRepoFake()
	lead := &domain.Lead{LeadID: "lead-1", Ruta: domain.RutaAsesor, Afiliado: true, Estado: domain.EstadoLeadCalificado}
	if err := repo.Crear(context.Background(), lead); err != nil {
		t.Fatal(err)
	}
	plan := &domain.PlanNutricion{PlanID: "plan-1", LeadID: lead.LeadID, Estado: domain.EstadoPlanActivo}
	out, err := (&ConsultarLeads{Leads: repo, Planes: planReadFake{plan: plan}}).Detalle(context.Background(), lead.LeadID)
	if err != nil || out.Semaforo != domain.SemaforoVerde || out.Plan == nil || out.Plan.PlanID != "plan-1" {
		t.Fatalf("detail=%+v err=%v", out, err)
	}
}

type fichaReadFake struct {
	ficha *domain.Ficha
	err   error
}

func (f fichaReadFake) Guardar(context.Context, *domain.Ficha) error { return nil }
func (f fichaReadFake) PorLead(context.Context, string) (*domain.Ficha, error) {
	return f.ficha, f.err
}
func TestConsultarLeadsQueueIsDeterministicAndReadOnly(t *testing.T) {
	repo := NuevoLeadRepoFake()
	now := time.Date(2026, 7, 26, 7, 0, 0, 0, time.UTC)
	leads := []*domain.Lead{
		{LeadID: "low", Nombre: "Low", Estado: domain.EstadoLeadCalificado, Ruta: domain.RutaAsesor, Prioridad: .20, ActualizadoEn: now},
		{LeadID: "high", Nombre: "High", Estado: domain.EstadoLeadCalificado, Ruta: domain.RutaAsesor, Afiliado: true, Prioridad: .90, Perfil: domain.Perfil{"categoria": {Valor: "A"}}, Capacidad: &domain.Capacidad{PresupuestoMax: 166800000}, Intencion: &domain.Intencion{Nivel: domain.NivelAlta}, ActualizadoEn: now},
		{LeadID: "nutrition", Ruta: domain.RutaNutricion, Prioridad: .80, ActualizadoEn: now},
	}
	for _, lead := range leads {
		if err := repo.Crear(context.Background(), lead); err != nil {
			t.Fatal(err)
		}
	}
	out, err := (&ConsultarLeads{Leads: repo}).Ejecutar(context.Background(), FiltroLeads{})
	if err != nil || len(out.Leads) != 3 {
		t.Fatalf("queue=%+v err=%v", out, err)
	}
	if out.Leads[0].LeadID != "high" || out.Leads[1].LeadID != "nutrition" || out.Leads[2].LeadID != "low" {
		t.Fatalf("order=%v", []string{out.Leads[0].LeadID, out.Leads[1].LeadID, out.Leads[2].LeadID})
	}
	if out.Cupo10.Usados != 1 || out.Cupo10.PorcentajeVentana != 10 || out.Leads[0].Semaforo != domain.SemaforoVerde {
		t.Fatalf("cupo/semaforo=%+v %+v", out.Cupo10, out.Leads[0])
	}
	if out.Leads[0].Resumen != "Afiliada cat. A · presupuesto $166.8M · intención alta" {
		t.Fatalf("summary=%q", out.Leads[0].Resumen)
	}
	stored, _ := repo.PorID(context.Background(), "high")
	if stored.Prioridad != .90 || stored.Estado != domain.EstadoLeadCalificado {
		t.Fatalf("read mutated lead=%+v", stored)
	}
}
func TestConsultarLeadsFiltersAndFichaDistinction(t *testing.T) {
	repo := NuevoLeadRepoFake()
	for _, lead := range []*domain.Lead{{LeadID: "a", Afiliado: true}, {LeadID: "b", Ruta: domain.RutaAsesor}} {
		if err := repo.Crear(context.Background(), lead); err != nil {
			t.Fatal(err)
		}
	}
	affiliate := true
	out, err := (&ConsultarLeads{Leads: repo}).Ejecutar(context.Background(), FiltroLeads{Afiliado: &affiliate})
	if err != nil || len(out.Leads) != 1 || out.Leads[0].LeadID != "a" {
		t.Fatalf("filter=%+v err=%v", out, err)
	}
	uc := &ConsultarLeads{Leads: repo, Fichas: fichaReadFake{err: &NotFoundError{Resource: "ficha", ID: "a"}}}
	if _, err := uc.Ficha(context.Background(), "a"); !errors.As(err, new(*NotFoundError)) {
		t.Fatalf("missing ficha=%v", err)
	}
	if _, err := uc.Ficha(context.Background(), "missing"); !errors.As(err, new(*NotFoundError)) {
		t.Fatalf("missing lead=%v", err)
	}
}
