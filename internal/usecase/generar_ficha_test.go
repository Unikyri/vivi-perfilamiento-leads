package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

type fichaRepoFake struct {
	stored          *domain.Ficha
	porErr, saveErr error
	calls           []string
}

func (f *fichaRepoFake) Guardar(_ context.Context, ficha *domain.Ficha) error {
	f.calls = append(f.calls, "ficha-save")
	if f.saveErr != nil {
		return f.saveErr
	}
	copy := *ficha
	f.stored = &copy
	return nil
}
func (f *fichaRepoFake) PorLead(_ context.Context, _ string) (*domain.Ficha, error) {
	f.calls = append(f.calls, "ficha-read")
	if f.porErr != nil {
		return nil, f.porErr
	}
	if f.stored == nil {
		return nil, &NotFoundError{Resource: "ficha", ID: "lead"}
	}
	return f.stored, nil
}

func fichaLead(id string, affiliate bool) *domain.Lead {
	lead := calLead(id, affiliate, domain.NivelAlta, 120)
	lead.Ruta = domain.RutaAsesor
	lead.Nombre, lead.Telefono = "Ana", "3001234567"
	lead.Perfil = calProfile(120, affiliate, "", "", false)
	lead.Perfil["ingreso_hogar"] = domain.CampoPerfil{Valor: int64(3000000), Fuente: domain.FuenteCampoVerificadoBase}
	lead.Perfil["categoria"] = domain.CampoPerfil{Valor: "B", Fuente: domain.FuenteCampoDeclarado}
	lead.Perfil["arriendo_actual"] = domain.CampoPerfil{Valor: int64(1200000), Fuente: domain.FuenteCampoDeclarado}
	lead.Intencion = &domain.Intencion{Nivel: domain.NivelAlta, Confianza: domain.NivelAlta, Senales: []string{"compra"}}
	lead.ConsumeCupo10 = !affiliate
	return lead
}
func fichaCatalog(rateDesistidos int) *calCatalogFake {
	catalog := calCatalog(100)
	for i := 0; i < 5; i++ {
		catalog.buyers = append(catalog.buyers, domain.Comprador{ID: i + 1, ProyectoID: "p1", PersonasACargo: 1, Desistio: i < rateDesistidos})
	}
	return catalog
}
func runFicha(t *testing.T, ctx context.Context, lead *domain.Lead, catalog *calCatalogFake, fichas *fichaRepoFake, saveErr error) (*domain.Ficha, *calRepoFake, []string, error) {
	t.Helper()
	raw := NuevoLeadRepoFake()
	if err := raw.Crear(context.Background(), lead); err != nil {
		t.Fatal(err)
	}
	order := []string{}
	leads := &calRepoFake{LeadRepoFake: raw, saveErr: saveErr, order: &order}
	uc := &GenerarFicha{Leads: leads, Fichas: fichas, Catalogo: catalog, IDs: NuevoIDFake("ficha"), Reloj: NuevoRelojFake(time.Date(2026, 7, 25, 12, 0, 0, 0, time.UTC))}
	out, err := uc.Ejecutar(ctx, lead.LeadID)
	return out, leads, order, err
}

func TestGenerarFichaGuardsAndReadFailuresDoNotWrite(t *testing.T) {
	cases := []struct {
		name   string
		lead   func() *domain.Lead
		cat    *calCatalogFake
		ferr   error
		cancel bool
	}{
		{"blank id", func() *domain.Lead { return fichaLead("", false) }, fichaCatalog(0), nil, false},
		{"wrong state", func() *domain.Lead { l := fichaLead("l", false); l.Estado = domain.EstadoLeadPerfilando; return l }, fichaCatalog(0), nil, false},
		{"wrong route", func() *domain.Lead { l := fichaLead("l", false); l.Ruta = domain.RutaNutricion; return l }, fichaCatalog(0), nil, false},
		{"missing intention", func() *domain.Lead { l := fichaLead("l", false); l.Intencion = nil; return l }, fichaCatalog(0), nil, false},
		{"missing capacity", func() *domain.Lead { l := fichaLead("l", false); l.Capacidad = nil; return l }, fichaCatalog(0), nil, false},
		{"ficha read failure", func() *domain.Lead { return fichaLead("l", false) }, fichaCatalog(0), errors.New("ficha unavailable"), false},
		{"projects failure", func() *domain.Lead { return fichaLead("l", false) }, &calCatalogFake{projectErr: errors.New("projects")}, nil, false},
		{"buyers failure", func() *domain.Lead { return fichaLead("l", false) }, &calCatalogFake{projects: calCatalog(0).projects, buyerErr: errors.New("buyers")}, nil, false},
		{"ficha save failure", func() *domain.Lead { return fichaLead("l", false) }, fichaCatalog(0), errors.New("ficha save"), false},
		{"canceled", func() *domain.Lead { return fichaLead("l", false) }, fichaCatalog(0), nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.cancel {
				cancel()
			}
			ferr, saveErr := tc.ferr, error(nil)
			if tc.name == "ficha save failure" {
				saveErr, ferr = ferr, nil
			}
			fichas := &fichaRepoFake{porErr: ferr, saveErr: saveErr}
			_, _, order, err := runFicha(t, ctx, tc.lead(), tc.cat, fichas, nil)
			if err == nil {
				t.Fatal("expected error")
			}
			if fichas.stored != nil || len(order) != 0 {
				t.Fatalf("writes = ficha:%v calls:%v lead:%v", fichas.stored, fichas.calls, order)
			}
		})
	}
}

func TestGenerarFichaContentParityThresholdAndNoAliasing(t *testing.T) {
	lead := fichaLead("l", true)
	catalog := fichaCatalog(1)
	fichas := &fichaRepoFake{}
	got, leads, order, err := runFicha(t, context.Background(), lead, catalog, fichas, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(order, []string{"save"}) {
		t.Fatalf("ordering = %v", order)
	}
	if got.FichaID != "ficha-1" || got.LeadID != "l" || got.Identificacion.Nombre != "Ana" || got.Identificacion.Categoria != "B" || !got.Identificacion.Afiliada {
		t.Fatalf("identity = %+v", got.Identificacion)
	}
	if got.BandaAdvertencia != nil {
		t.Fatalf("unexpected warning: %v", *got.BandaAdvertencia)
	}
	if !reflect.DeepEqual(got.Beneficios, []string{"Subsidio de caja 30 SMMLV", "Crédito propio Colsubsidio", "Acompañamiento PerteneSer"}) {
		t.Fatalf("benefits = %#v", got.Beneficios)
	}
	if !reflect.DeepEqual(got.ArgumentosVenta, []string{"Paga $1.200.000 de arriendo; la cuota estimada es $1.200.000"}) {
		t.Fatalf("arguments = %#v", got.ArgumentosVenta)
	}
	if got.AlertaDesistimiento.TasaVecinos != .2 || got.AlertaDesistimiento.Activa || got.AlertaDesistimiento.Detalle != nil {
		t.Fatalf("threshold = %+v", got.AlertaDesistimiento)
	}
	want := construirDecision(lead, catalog.projects, catalog.buyers)
	if !reflect.DeepEqual(got.Recomendaciones, want.recomendaciones) {
		t.Fatalf("recommendation parity = %#v/%#v", got.Recomendaciones, want.recomendaciones)
	}
	got.Perfil["nested"] = domain.CampoPerfil{Valor: map[string]any{"x": []int{1}}}
	got.Perfil["nested"].Valor.(map[string]any)["x"].([]int)[0] = 9
	got.Intencion.Senales[0] = "changed"
	stored, _ := leads.PorID(context.Background(), "l")
	if stored.Perfil["nested"].Valor != nil || stored.Intencion.Senales[0] != "compra" {
		t.Fatal("output aliases lead state")
	}

	low := fichaLead("low", false)
	for key, field := range low.Perfil {
		field.Fuente = domain.FuenteCampoDeclarado
		low.Perfil[key] = field
	}
	low.Perfil["arriendo_actual"] = domain.CampoPerfil{Valor: int64(0), Fuente: domain.FuenteCampoDeclarado}
	lowFichas := &fichaRepoFake{}
	lowOut, _, _, err := runFicha(t, context.Background(), low, fichaCatalog(2), lowFichas, nil)
	if err != nil {
		t.Fatal(err)
	}
	if lowOut.BandaAdvertencia == nil || *lowOut.BandaAdvertencia != advertenciaPerfilParcial || len(lowOut.ArgumentosVenta) != 0 || !lowOut.AlertaDesistimiento.Activa {
		t.Fatalf("low/rent/alert = %+v/%#v/%+v", lowOut.BandaAdvertencia, lowOut.ArgumentosVenta, lowOut.AlertaDesistimiento)
	}
}

func TestGenerarFichaFichaFirstFailureAndRetry(t *testing.T) {
	lead := fichaLead("retry", false)
	fichas := &fichaRepoFake{}
	first, leads, order, err := runFicha(t, context.Background(), lead, fichaCatalog(0), fichas, errors.New("cas"))
	if err == nil || fichas.stored == nil || first.FichaID != "ficha-1" || len(order) != 1 || order[0] != "save" {
		t.Fatalf("first attempt = %#v/%v/%v/%v", first.FichaID, err, fichas.stored, order)
	}
	storedLead, _ := leads.PorID(context.Background(), "retry")
	if storedLead.Estado != domain.EstadoLeadCalificado || storedLead.Ruta != domain.RutaAsesor {
		t.Fatalf("lead changed after CAS failure: %+v", storedLead)
	}
	leads.saveErr = nil
	uc := &GenerarFicha{Leads: leads, Fichas: fichas, Catalogo: fichaCatalog(0), IDs: NuevoIDFake("new"), Reloj: NuevoRelojFake(time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC))}
	second, err := uc.Ejecutar(context.Background(), "retry")
	if err != nil || second.FichaID != first.FichaID || !second.GeneradaEn.Equal(first.GeneradaEn) {
		t.Fatalf("retry identity/time = %#v/%v", second, err)
	}
	storedLead, _ = leads.PorID(context.Background(), "retry")
	wantUpdated := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	if storedLead.Estado != domain.EstadoLeadEntregado || !storedLead.ActualizadoEn.Equal(wantUpdated) || len(fichas.calls) != 4 || !reflect.DeepEqual(*leads.order, []string{"save", "save"}) {
		t.Fatalf("retry ordering/state/timestamp = %v/%v/%v/%v", fichas.calls, *leads.order, storedLead.Estado, storedLead.ActualizadoEn)
	}
}
