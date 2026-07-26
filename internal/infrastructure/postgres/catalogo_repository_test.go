package postgres

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type countingFS struct {
	fsys  fs.FS
	reads int
}

func (f *countingFS) Open(name string) (fs.File, error) {
	f.reads++
	return f.fsys.Open(name)
}

func writeCatalogFiles(t *testing.T, root string) {
	t.Helper()
	files := map[string]string{
		"proyectos.json":      `[{"proyecto_id":"p-1","nombre":"Uno"}]`,
		"compradores.json":    `[{"id":1,"proyecto_id":"p-1","proyecto":"Uno"}]`,
		"afiliados_mock.json": `[{"cedula":"1","nombre":"Ana"}]`,
		"brochures/p-1.md":    "---\nproyecto_id: p-1\n---\n# Uno\n",
	}
	for name, content := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCatalogoRepository_EagerSnapshotAndDefensiveResults(t *testing.T) {
	root := t.TempDir()
	writeCatalogFiles(t, root)
	counted := &countingFS{fsys: os.DirFS(root)}
	repository, err := NuevoCatalogo(counted)
	if err != nil {
		t.Fatal(err)
	}
	reads := counted.reads
	projects, _ := repository.Proyectos(context.Background())
	buyers, _ := repository.Compradores(context.Background())
	affiliate, _ := repository.AfiliadoPorCedula(context.Background(), "1")
	if _, err := repository.BrochureMarkdown(context.Background(), "p-1"); err != nil {
		t.Fatal(err)
	}
	projects["p-1"] = domain.Proyecto{}
	buyers[0].Proyecto = "mutated"
	affiliate.Nombre = "mutated"
	again, _ := repository.Proyectos(context.Background())
	buyersAgain, _ := repository.Compradores(context.Background())
	affiliateAgain, _ := repository.AfiliadoPorCedula(context.Background(), "1")
	if counted.reads != reads || again["p-1"].Nombre != "Uno" || buyersAgain[0].Proyecto != "Uno" || affiliateAgain.Nombre != "Ana" {
		t.Fatalf("snapshot was reread or mutated: reads %d/%d, %#v, %#v, %#v", counted.reads, reads, again, buyersAgain, affiliateAgain)
	}
}

func TestCatalogoRepository_SourceFailureAndTypedMisses(t *testing.T) {
	if _, err := NuevoCatalogo(os.DirFS(t.TempDir())); err == nil {
		t.Fatal("missing source files must fail initialization")
	}
	// A comprador referencing a proyecto_id absent from the catalog is real
	// historical data (~20% of the production dataset), not corruption: it
	// must NOT fail initialization, and must still be readable — GemeloKNN
	// treats it as a valid kNN comparison point with no resolvable zone.
	orphanRoot := t.TempDir()
	writeCatalogFiles(t, orphanRoot)
	if err := os.WriteFile(filepath.Join(orphanRoot, "compradores.json"), []byte(`[{"id":1,"proyecto_id":"missing"}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	orphanRepo, err := NuevoCatalogo(os.DirFS(orphanRoot))
	if err != nil {
		t.Fatalf("buyer with unknown project must not fail initialization: %v", err)
	}
	orphanBuyers, err := orphanRepo.Compradores(context.Background())
	if err != nil || len(orphanBuyers) != 1 || orphanBuyers[0].ProyectoID != "missing" {
		t.Fatalf("orphan comprador must still be retrievable: buyers=%#v err=%v", orphanBuyers, err)
	}

	root := t.TempDir()
	writeCatalogFiles(t, root)
	repository, err := NuevoCatalogo(os.DirFS(root))
	if err != nil {
		t.Fatal(err)
	}
	for name, call := range map[string]func() error{
		"affiliate": func() error { _, err := repository.AfiliadoPorCedula(context.Background(), "missing"); return err },
		"brochure":  func() error { _, err := repository.BrochureMarkdown(context.Background(), "missing"); return err },
	} {
		err := call()
		var notFound *usecase.NotFoundError
		if !errors.As(err, &notFound) || !errors.Is(err, usecase.ErrNoEncontrado) || notFound.Resource == "" {
			t.Errorf("%s miss = %v", name, err)
		}
	}
}
