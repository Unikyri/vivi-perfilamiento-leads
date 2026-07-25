package usecase

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

func testLead(id string, priority float64, affiliate bool, route domain.Ruta) *domain.Lead {
	return &domain.Lead{LeadID: id, Prioridad: priority, Afiliado: affiliate, Ruta: route, Perfil: domain.Perfil{}}
}
func TestLeadRepoFake_LifecycleCASAndDeepIsolation(t *testing.T) {
	ctx, repo := context.Background(), NuevoLeadRepoFake()
	n := int64(7)
	lead := testLead("a", 1, true, domain.RutaAsesor)
	lead.Version = 1
	lead.Perfil["nested"] = domain.CampoPerfil{Valor: map[string]any{"nums": []int{1}, "ptr": &n, "amount": int64(4)}}
	lead.Capacidad = &domain.Capacidad{Desglose: []domain.ItemDesglose{{Monto: 8}}}
	lead.Intencion = &domain.Intencion{Senales: []string{"first"}}
	if err := repo.Crear(ctx, lead); err != nil || lead.Version != 1 {
		t.Fatalf("Crear() error/version = %v/%d", err, lead.Version)
	}
	if err := repo.Crear(ctx, lead); !errors.Is(err, errLeadDuplicado) {
		t.Fatalf("duplicate error = %v", err)
	}
	for _, version := range []int{2, -1} {
		if err := repo.Crear(ctx, &domain.Lead{LeadID: "bad", Version: version}); !errors.Is(err, errVersionInicial) {
			t.Errorf("version %d error = %v", version, err)
		}
	}
	m := lead.Perfil["nested"].Valor.(map[string]any)
	m["nums"].([]int)[0] = 9
	n = 99
	lead.Capacidad.Desglose[0].Monto = 99
	lead.Intencion.Senales[0] = "changed"
	got, err := repo.PorID(ctx, "a")
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != 1 || m["nums"].([]int)[0] != 9 || got.Perfil["nested"].Valor.(map[string]any)["nums"].([]int)[0] != 1 || got.Perfil["nested"].Valor.(map[string]any)["ptr"].(*int64) == &n || got.Perfil["nested"].Valor.(map[string]any)["amount"].(int64) != 4 || got.Capacidad.Desglose[0].Monto != 8 || got.Intencion.Senales[0] != "first" {
		t.Fatalf("input alias leaked: %#v", got)
	}
	got.Perfil["nested"].Valor.(map[string]any)["nums"].([]int)[0] = 3
	got.Capacidad.Desglose[0].Monto = 4
	fresh, _ := repo.PorID(ctx, "a")
	if fresh.Perfil["nested"].Valor.(map[string]any)["nums"].([]int)[0] != 1 || fresh.Capacidad.Desglose[0].Monto != 8 {
		t.Fatal("output alias leaked")
	}
	stale := cloneLead(fresh)
	current := cloneLead(fresh)
	current.Nombre = "saved"
	if err := repo.Guardar(ctx, current); err != nil || current.Version != 2 {
		t.Fatalf("Guardar() = %v, version %d", err, current.Version)
	}
	if err := repo.Guardar(ctx, stale); !errors.Is(err, ErrOptimisticLock) || stale.Version != 1 {
		t.Fatalf("stale CAS = %v, version %d", err, stale.Version)
	}
	final, _ := repo.PorID(ctx, "a")
	if final.Nombre != "saved" || final.Version != 2 {
		t.Fatalf("stale write changed storage: %#v", final)
	}
	missing, err := repo.PorID(ctx, "missing")
	if missing != nil || !errors.Is(err, ErrNoEncontrado) {
		t.Fatalf("absence = %v/%v", missing, err)
	}
}
func TestLeadRepoFake_ListarFiltersOrderAndEmpty(t *testing.T) {
	repo, ctx := NuevoLeadRepoFake(), context.Background()
	for _, lead := range []*domain.Lead{testLead("b", 5, true, domain.RutaAsesor), testLead("a", 5, true, domain.RutaAsesor), testLead("c", 9, false, domain.RutaNutricion)} {
		if err := repo.Crear(ctx, lead); err != nil || lead.Version != 1 {
			t.Fatalf("Crear() error/version = %v/%d", err, lead.Version)
		}
	}
	yes, route := true, domain.RutaAsesor
	got, _ := repo.Listar(ctx, FiltroLeads{Afiliado: &yes, Ruta: &route})
	if len(got) != 2 || got[0].LeadID != "a" || got[1].LeadID != "b" {
		t.Fatalf("filter/order = %#v", got)
	}
	no := false
	empty, _ := repo.Listar(ctx, FiltroLeads{Afiliado: &no, Ruta: &route})
	if empty == nil || len(empty) != 0 {
		t.Fatalf("empty = %#v", empty)
	}
}
func TestLeadRepoFake_ConversationAndMinimalFakes(t *testing.T) {
	ctx, repo := context.Background(), NuevoLeadRepoFake()
	if err := repo.Crear(ctx, testLead("a", 0, false, domain.RutaAsesor)); err != nil {
		t.Fatal(err)
	}
	attachment := map[string]any{"items": []string{"x"}}
	for _, message := range []*domain.Mensaje{{LeadID: "a", MensajeID: "late", CreadoEn: time.Unix(2, 0), Adjunto: attachment}, {LeadID: "a", MensajeID: "early", CreadoEn: time.Unix(1, 0)}} {
		if err := repo.AgregarMensaje(ctx, message); err != nil {
			t.Fatal(err)
		}
	}
	attachment["items"].([]string)[0] = "input-changed"
	conversation, _ := repo.Conversacion(ctx, "a")
	if len(conversation) != 2 || conversation[0].MensajeID != "early" || conversation[1].MensajeID != "late" {
		t.Fatalf("conversation order = %#v", conversation)
	}
	conversation[1].Adjunto["items"].([]string)[0] = "changed"
	reread, _ := repo.Conversacion(ctx, "a")
	if reread[1].Adjunto["items"].([]string)[0] != "x" {
		t.Fatal("message output alias leaked")
	}
}
func TestLeadRepoFake_ConcurrentAccess(t *testing.T) {
	ctx, repo := context.Background(), NuevoLeadRepoFake()
	_ = repo.Crear(ctx, testLead("a", 1, false, domain.RutaAsesor))
	var group sync.WaitGroup
	for i := 0; i < 16; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			lead, _ := repo.PorID(ctx, "a")
			lead.Nombre = "updated"
			_ = repo.Guardar(ctx, lead)
			_, _ = repo.Listar(ctx, FiltroLeads{})
		}()
	}
	group.Wait()
}
