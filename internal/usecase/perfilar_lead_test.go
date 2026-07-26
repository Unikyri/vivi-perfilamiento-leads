package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

type catalogFake struct {
	affiliates map[string]*Afiliado
	err        error
}

func (f *catalogFake) Proyectos(context.Context) (map[string]domain.Proyecto, error) { return nil, nil }
func (f *catalogFake) Compradores(context.Context) ([]domain.Comprador, error)       { return nil, nil }
func (f *catalogFake) BrochureMarkdown(context.Context, string) (string, error)      { return "", nil }
func (f *catalogFake) AfiliadoPorCedula(_ context.Context, cedula string) (*Afiliado, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.affiliates[cedula], nil
}

type busFake struct{ events []Evento }

func (b *busFake) Publicar(_ context.Context, event Evento)        { b.events = append(b.events, event) }
func (b *busFake) Suscribir(string, func(context.Context, Evento)) {}

func ana() *Afiliado {
	return &Afiliado{Cedula: "1032456789", Nombre: "Ana", Categoria: "B", Segmento: "A", IngresoMensual: 2_600_000, PersonasACargo: 2, TipoHogar: "FAMILIAR", AfiliadoActivo: true}
}
func newUC(c CatalogoRepository, r LeadRepository, b BusEventos) *PerfilarLead {
	return &PerfilarLead{Leads: r, Catalogo: c, IDs: NuevoIDFake("lead"), Bus: b, Reloj: NuevoRelojFake(time.Unix(100, 0))}
}

func TestPerfilarLeadEjecucion(t *testing.T) {
	cases := []struct {
		name       string
		catalog    CatalogoRepository
		cedula     string
		affiliated bool
		profile    int
		subsidy    int64
	}{
		{"active affiliate", &catalogFake{affiliates: map[string]*Afiliado{ana().Cedula: ana()}}, "1032456789", true, 7, 52_527_150},
		{"unknown fallback", &catalogFake{affiliates: map[string]*Afiliado{}}, "000", false, 0, 0},
		{"inactive fallback", &catalogFake{affiliates: map[string]*Afiliado{"x": {Cedula: "x", AfiliadoActivo: false}}}, "x", false, 0, 0},
		{"catalog failure fallback", &catalogFake{err: errors.New("catalog down")}, "x", false, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo, bus := NuevoLeadRepoFake(), &busFake{}
			out, err := newUC(tc.catalog, repo, bus).Ejecutar(context.Background(), EntradaPerfilar{Cedula: tc.cedula})
			if err != nil {
				t.Fatal(err)
			}
			lead, err := repo.PorID(context.Background(), out.LeadID)
			if err != nil || lead.Estado != domain.EstadoLeadPerfilando || lead.Afiliado != tc.affiliated {
				t.Fatalf("lead=%+v err=%v", lead, err)
			}
			if len(lead.Perfil) != tc.profile || lead.Capacidad.SubsidioAplicable != tc.subsidy || lead.Version != 1 {
				t.Fatalf("profile=%v capacity=%+v version=%d", lead.Perfil, lead.Capacidad, lead.Version)
			}
			if tc.affiliated {
				people, _ := lead.Perfil.Entero(campoPersonas)
				owns, _ := lead.Perfil.Booleano(campoVivienda)
				prior, _ := lead.Perfil.Booleano(campoSubsidio)
				if lead.Nombre != "Ana" || people != 3 || owns || prior || !lead.Perfil.EsVerificado(campoIngreso) || !lead.Perfil.EsVerificado(campoCategoria) || !lead.Perfil.EsVerificado(campoSegmento) || !lead.Perfil.EsVerificado(campoPersonas) || !lead.Perfil.EsVerificado(campoTipoHogar) || !lead.Perfil.EsVerificado(campoVivienda) || !lead.Perfil.EsVerificado(campoSubsidio) {
					t.Fatalf("unexpected verified pre-profile: %+v", lead)
				}
			}
			wantEvent := map[string]any{"cedula": lead.Cedula, "nombre": lead.Nombre, "telefono": lead.Telefono, "fuente": lead.Fuente}
			if len(bus.events) != 1 || bus.events[0].Tipo != EvLeadNuevo || bus.events[0].LeadID != out.LeadID || !reflect.DeepEqual(bus.events[0].Payload, wantEvent) {
				t.Fatalf("events=%+v", bus.events)
			}
		})
	}
}

func TestPerfilarLeadErrorsDoNotPublish(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	repo, bus := NuevoLeadRepoFake(), &busFake{}
	if _, err := newUC(&catalogFake{affiliates: map[string]*Afiliado{}}, repo, bus).Ejecutar(ctx, EntradaPerfilar{}); !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
	bad := &errorRepo{LeadRepoFake: NuevoLeadRepoFake(), createErr: errors.New("create")}
	if _, err := newUC(&catalogFake{}, bad, bus).Ejecutar(context.Background(), EntradaPerfilar{}); err == nil || len(bus.events) != 0 {
		t.Fatalf("err=%v events=%d", err, len(bus.events))
	}
}

type errorRepo struct {
	*LeadRepoFake
	createErr, saveErr, porIDErr error
}

func (r *errorRepo) Crear(context.Context, *domain.Lead) error { return r.createErr }
func (r *errorRepo) PorID(ctx context.Context, id string) (*domain.Lead, error) {
	if r.porIDErr != nil {
		return nil, r.porIDErr
	}
	return r.LeadRepoFake.PorID(ctx, id)
}
func (r *errorRepo) Guardar(ctx context.Context, lead *domain.Lead) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	return r.LeadRepoFake.Guardar(ctx, lead)
}

func TestReconsultarFamiliar(t *testing.T) {
	a := ana()
	family := &Afiliado{Cedula: "1015789456", IngresoMensual: 1_900_000, AfiliadoActivo: true}
	cat := &catalogFake{affiliates: map[string]*Afiliado{a.Cedula: a, family.Cedula: family}}
	repo, bus := NuevoLeadRepoFake(), &busFake{}
	uc := newUC(cat, repo, bus)
	out, err := uc.Ejecutar(context.Background(), EntradaPerfilar{Cedula: a.Cedula})
	if err != nil {
		t.Fatal(err)
	}
	if err = uc.ReconsultarPorFamiliar(context.Background(), out.LeadID, family.Cedula); err != nil {
		t.Fatal(err)
	}
	lead, _ := repo.PorID(context.Background(), out.LeadID)
	income, _ := lead.Perfil.Entero(campoIngreso)
	version := lead.Version
	hogar, _ := lead.Perfil.Booleano(campoHogar)
	if income != 4_500_000 || !hogar || !lead.Perfil.EsVerificado(campoHogar) || !lead.Perfil.EsVerificado(campoFamiliar) || lead.Capacidad.SubsidioAplicable != 35_018_100 || len(bus.events) != 1 {
		t.Fatalf("lead=%+v events=%d", lead, len(bus.events))
	}
	if err = uc.ReconsultarPorFamiliar(context.Background(), out.LeadID, family.Cedula); err != nil {
		t.Fatal(err)
	}
	repeat, _ := repo.PorID(context.Background(), out.LeadID)
	repeatIncome, _ := repeat.Perfil.Entero(campoIngreso)
	if repeatIncome != income || repeat.Version != version || repeat.Capacidad.SubsidioAplicable != lead.Capacidad.SubsidioAplicable || len(bus.events) != 1 {
		t.Fatalf("repeat income=%d version=%d capacity=%+v events=%d", repeatIncome, repeat.Version, repeat.Capacidad, len(bus.events))
	}
	if err = uc.ReconsultarPorFamiliar(context.Background(), out.LeadID, "other"); !errors.Is(err, ErrFamiliarYaRegistrado) {
		t.Fatalf("err=%v", err)
	}
}

func TestReconsultarFamiliarNoEncontradoPersisteDeclarado(t *testing.T) {
	repo, bus := NuevoLeadRepoFake(), &busFake{}
	uc := newUC(&catalogFake{}, repo, bus)
	out, err := uc.Ejecutar(context.Background(), EntradaPerfilar{})
	if err != nil {
		t.Fatal(err)
	}
	err = uc.ReconsultarPorFamiliar(context.Background(), out.LeadID, "missing")
	if !errors.Is(err, ErrFamiliarNoEncontrado) {
		t.Fatalf("err=%v", err)
	}
	lead, _ := repo.PorID(context.Background(), out.LeadID)
	field := lead.Perfil[campoFamiliar]
	if field.Fuente != domain.FuenteCampoDeclarado || !field.RequiereConfirmacion || len(bus.events) != 1 {
		t.Fatalf("field=%+v events=%d", field, len(bus.events))
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err = uc.ReconsultarPorFamiliar(ctx, out.LeadID, "again"); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled err=%v", err)
	}
	loading := &errorRepo{LeadRepoFake: repo, porIDErr: errors.New("load")}
	if err = newUC(&catalogFake{}, loading, bus).ReconsultarPorFamiliar(context.Background(), out.LeadID, "again"); err == nil {
		t.Fatal("expected PorID error")
	}
	failing := &errorRepo{LeadRepoFake: repo, saveErr: errors.New("cas")}
	if err = newUC(&catalogFake{}, failing, bus).ReconsultarPorFamiliar(context.Background(), out.LeadID, "again"); err == nil || errors.Is(err, ErrFamiliarNoEncontrado) {
		t.Fatalf("err=%v", err)
	}
}
