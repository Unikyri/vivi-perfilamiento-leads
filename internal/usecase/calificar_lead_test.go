package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain/motor"
)

type calCatalogFake struct {
	catalogoRepositoryShape
	projects             map[string]domain.Proyecto
	buyers               []domain.Comprador
	projectErr, buyerErr error
}

func (f *calCatalogFake) Proyectos(context.Context) (map[string]domain.Proyecto, error) {
	return f.projects, f.projectErr
}
func (f *calCatalogFake) Compradores(context.Context) ([]domain.Comprador, error) {
	return f.buyers, f.buyerErr
}

type calBusFake struct {
	events []Evento
	order  *[]string
}

func (b *calBusFake) Publicar(_ context.Context, event Evento) {
	b.events = append(b.events, event)
	if b.order != nil {
		*b.order = append(*b.order, "event")
	}
}
func (b *calBusFake) Suscribir(string, func(context.Context, Evento)) {}

type calRepoFake struct {
	*LeadRepoFake
	saveErr error
	order   *[]string
}

func (r *calRepoFake) Guardar(ctx context.Context, lead *domain.Lead) error {
	if r.order != nil {
		*r.order = append(*r.order, "save")
	}
	if r.saveErr != nil {
		return r.saveErr
	}
	return r.LeadRepoFake.Guardar(ctx, lead)
}
func calProfile(resources int64, affiliate bool, labor, caja string, hogar bool) domain.Perfil {
	return domain.Perfil{"ingreso_hogar": {Valor: int64(0), Fuente: domain.FuenteCampoVerificadoBase}, "recursos_propios": {Valor: resources, Fuente: domain.FuenteCampoVerificadoBase}, "tiene_vivienda": {Valor: false, Fuente: domain.FuenteCampoVerificadoBase}, "recibio_subsidio": {Valor: false, Fuente: domain.FuenteCampoVerificadoBase}, "personas_hogar": {Valor: int64(2), Fuente: domain.FuenteCampoVerificadoBase}, "edad": {Valor: int64(35), Fuente: domain.FuenteCampoVerificadoBase}, "zona_deseada": {Valor: "NORTE", Fuente: domain.FuenteCampoVerificadoBase}, "situacion_laboral": {Valor: labor}, "caja_externa": {Valor: caja}, "hogar_con_afiliado": {Valor: hogar}}
}
func calLead(id string, affiliate bool, level domain.Nivel, resources int64) *domain.Lead {
	return &domain.Lead{LeadID: id, Estado: domain.EstadoLeadCalificado, Afiliado: affiliate, Perfil: calProfile(resources, affiliate, "", "", false), Intencion: &domain.Intencion{Nivel: level}, Capacidad: &domain.Capacidad{}}
}
func calCatalog(price int64) *calCatalogFake {
	return &calCatalogFake{projects: map[string]domain.Proyecto{"p1": {ProyectoID: "p1", Nombre: "P1", Zona: "NORTE", PrecioDesde: price}}}
}
func runCal(t *testing.T, ctx context.Context, lead *domain.Lead, catalog CatalogoRepository, saveErr error) (SalidaCalificar, *calBusFake, []string, error) {
	raw := NuevoLeadRepoFake()
	_ = raw.Crear(context.Background(), lead)
	order := []string{}
	bus := &calBusFake{order: &order}
	uc := &CalificarLead{Leads: &calRepoFake{LeadRepoFake: raw, saveErr: saveErr, order: &order}, Catalogo: catalog, Bus: bus}
	out, err := uc.Ejecutar(ctx, EntradaCalificar{LeadID: lead.LeadID})
	return out, bus, order, err
}
func require(t *testing.T, ok bool, format string, args ...any) {
	if !ok {
		t.Fatalf(format, args...)
	}
}
func TestCalificarLeadGuardsAndReadFailuresDoNotWrite(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		mutate   func(*domain.Lead)
		catalog  *calCatalogFake
		canceled bool
	}{
		{name: "blank id", catalog: calCatalog(1)},
		{name: "wrong state", input: "l", mutate: func(l *domain.Lead) { l.Estado = domain.EstadoLeadPerfilando }, catalog: calCatalog(1)},
		{name: "missing intention", input: "l", mutate: func(l *domain.Lead) { l.Intencion = nil }, catalog: calCatalog(1)},
		{name: "missing capacity", input: "l", mutate: func(l *domain.Lead) { l.Capacidad = nil }, catalog: calCatalog(1)},
		{name: "projects failure", input: "l", catalog: &calCatalogFake{projectErr: errors.New("catalog")}},
		{name: "buyers failure", input: "l", catalog: &calCatalogFake{projects: calCatalog(1).projects, buyerErr: errors.New("buyers")}},
		{name: "canceled", input: "l", catalog: calCatalog(1), canceled: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lead := calLead("l", false, domain.NivelAlta, 100)
			if tc.mutate != nil {
				tc.mutate(lead)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.canceled {
				cancel()
			}
			lead.LeadID = tc.input
			_, bus, order, err := runCal(t, ctx, lead, tc.catalog, nil)
			require(t, err != nil, "expected error")
			require(t, len(order) == 0 && len(bus.events) == 0, "writes/events = %v/%d", order, len(bus.events))
		})
	}
}
func TestCalificarLeadCandidateKNNAndDeterminism(t *testing.T) {
	catalog := calCatalog(80)
	catalog.projects["expensive"] = domain.Proyecto{ProyectoID: "expensive", PrecioDesde: 90}
	for i := 0; i < 31; i++ {
		catalog.buyers = append(catalog.buyers, domain.Comprador{ID: i + 1, ProyectoID: "p1", PersonasACargo: 1})
	}
	out, _, _, err := runCal(t, context.Background(), calLead("l", false, domain.NivelAlta, 100), catalog, nil)
	require(t, err == nil, "qualification failed: %v", err)
	want := motor.CalcularCapacidad(calProfile(100, false, "", "", false), false, 80)
	require(t, out.Capacidad.Ratio == want.Ratio && len(out.Vecinos) == 30 && len(out.Recomendaciones) == 1, "candidate/knn = %+v/%d/%d", out.Capacidad, len(out.Vecinos), len(out.Recomendaciones))
	again := construirDecision(calLead("l", false, domain.NivelAlta, 100), catalog.projects, catalog.buyers)
	require(t, reflect.DeepEqual(out.Recomendaciones, again.recomendaciones), "recommendations not deterministic")
	fallback := construirDecision(calLead("l", false, domain.NivelAlta, 100), map[string]domain.Proyecto{"x": {ProyectoID: "x", PrecioDesde: 1000}}, nil)
	require(t, fallback.capacidad.Ratio == motor.CalcularCapacidad(calProfile(100, false, "", "", false), false, 0).Ratio, "zero candidate did not use median fallback")
}
func TestCalificarLeadRoutesConversionPriorityAndDurability(t *testing.T) {
	cases := []struct {
		name             string
		affiliate        bool
		level            domain.Nivel
		resources        int64
		labor, caja      string
		hogar            bool
		route            domain.Ruta
		state            domain.EstadoLead
		semaphore        domain.Semaforo
		cupo, conversion bool
		weight           float64
	}{
		{"affiliate advisor", true, domain.NivelAlta, 120, "", "", false, domain.RutaAsesor, domain.EstadoLeadCalificado, domain.SemaforoVerde, false, false, 1}, {"affiliate remarketing", true, domain.NivelBaja, 120, "", "", false, domain.RutaRemarketing, domain.EstadoLeadRemarketing, domain.SemaforoGris, false, false, .25},
		{"affiliate nutrition", true, domain.NivelAlta, 50, "", "", false, domain.RutaNutricion, domain.EstadoLeadEnNutricion, domain.SemaforoAmbar, false, false, .5}, {"affiliate goodbye", true, domain.NivelBaja, 50, "", "", false, domain.RutaDespedida, domain.EstadoLeadDespedido, domain.SemaforoGris, false, false, .1},
		{"independent converts", false, domain.NivelAlta, 120, "INDEPENDIENTE", "", false, domain.RutaNutricion, domain.EstadoLeadEnNutricion, domain.SemaforoAmbar, false, true, .5}, {"hogar converts", false, domain.NivelAlta, 120, "", "", true, domain.RutaNutricion, domain.EstadoLeadEnNutricion, domain.SemaforoAmbar, false, true, .5},
		{"caja converts", false, domain.NivelAlta, 120, "", "externa", false, domain.RutaNutricion, domain.EstadoLeadEnNutricion, domain.SemaforoAmbar, false, true, .5}, {"blank caja", false, domain.NivelAlta, 120, "", "   ", false, domain.RutaAsesor, domain.EstadoLeadCalificado, domain.SemaforoVerde, true, false, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lead := calLead("l", tc.affiliate, tc.level, tc.resources)
			lead.Perfil = calProfile(tc.resources, tc.affiliate, tc.labor, tc.caja, tc.hogar)
			out, bus, order, err := runCal(t, context.Background(), lead, calCatalog(100), nil)
			require(t, err == nil, "qualification failed: %v", err)
			wantPriority := tc.weight * min(out.Capacidad.Ratio, 1.2)
			require(t, out.Ruta == tc.route && out.Estado == tc.state && out.Semaforo == tc.semaphore && out.ConsumeCupo10 == tc.cupo && out.Conversion == tc.conversion && out.Prioridad == wantPriority, "out = %+v, want priority = %v", out, wantPriority)
			require(t, reflect.DeepEqual(order, []string{"save", "event"}) && len(bus.events) == 1, "ordering = %v, events = %d", order, len(bus.events))
			event := bus.events[0]
			require(t, event.Payload["ruta"] == tc.route && event.Payload["semaforo"] == tc.semaphore && event.Payload["consume_cupo_10"] == tc.cupo, "payload = %v", event.Payload)
		})
	}
	lead := calLead("fail", false, domain.NivelAlta, 120)
	_, bus, order, err := runCal(t, context.Background(), lead, calCatalog(100), errors.New("cas"))
	require(t, err != nil && len(bus.events) == 0 && reflect.DeepEqual(order, []string{"save"}), "save failure = %v/%v/%d", err, order, len(bus.events))
}
