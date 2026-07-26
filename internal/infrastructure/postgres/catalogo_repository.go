package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

type catalogSnapshot struct {
	proyectos   map[string]domain.Proyecto
	compradores []domain.Comprador
	afiliados   map[string]usecase.Afiliado
	brochures   map[string]string
}

// CatalogoRepository is an eager, read-only snapshot of the Contract §4 files.
type CatalogoRepository struct{ snapshot catalogSnapshot }

var _ usecase.CatalogoRepository = (*CatalogoRepository)(nil)

// NuevoCatalogo loads every catalog asset before returning a usable repository.
func NuevoCatalogo(data fs.FS) (*CatalogoRepository, error) {
	if data == nil {
		return nil, fmt.Errorf("catalogo: filesystem is nil")
	}
	var proyectos []domain.Proyecto
	if err := readJSON(data, "proyectos.json", &proyectos); err != nil {
		return nil, fmt.Errorf("catalogo: proyectos: %w", err)
	}
	var compradores []domain.Comprador
	if err := readJSON(data, "compradores.json", &compradores); err != nil {
		return nil, fmt.Errorf("catalogo: compradores: %w", err)
	}
	var afiliados []usecase.Afiliado
	if err := readJSON(data, "afiliados_mock.json", &afiliados); err != nil {
		return nil, fmt.Errorf("catalogo: afiliados: %w", err)
	}

	snapshot := catalogSnapshot{
		proyectos:   make(map[string]domain.Proyecto, len(proyectos)),
		compradores: append([]domain.Comprador(nil), compradores...),
		afiliados:   make(map[string]usecase.Afiliado, len(afiliados)),
		brochures:   make(map[string]string),
	}
	for _, proyecto := range proyectos {
		if proyecto.ProyectoID == "" {
			return nil, fmt.Errorf("catalogo: proyecto has empty proyecto_id")
		}
		if _, exists := snapshot.proyectos[proyecto.ProyectoID]; exists {
			return nil, fmt.Errorf("catalogo: duplicate proyecto_id %q", proyecto.ProyectoID)
		}
		snapshot.proyectos[proyecto.ProyectoID] = proyecto
	}
	buyerIDs := make(map[int]struct{}, len(compradores))
	huerfanos := make(map[string]int)
	for _, comprador := range compradores {
		if comprador.ID == 0 {
			return nil, fmt.Errorf("catalogo: comprador has empty id")
		}
		if _, exists := buyerIDs[comprador.ID]; exists {
			return nil, fmt.Errorf("catalogo: duplicate comprador id %d", comprador.ID)
		}
		// A comprador whose proyecto_id has no matching catalog entry is
		// legitimate historical data (a real transaction outside the current
		// 16-project marketing catalog), not corruption: GemeloKNN and
		// RecomendarProyectos already treat these buyers as valid kNN
		// comparison points that simply resolve to no catalog zone. Rejecting
		// them here would crash the app on boot against the real dataset,
		// where ~20% of compradores.json falls in this category.
		if _, exists := snapshot.proyectos[comprador.ProyectoID]; !exists {
			huerfanos[comprador.ProyectoID]++
		}
		buyerIDs[comprador.ID] = struct{}{}
	}
	if len(huerfanos) > 0 {
		ids := make([]string, 0, len(huerfanos))
		for id := range huerfanos {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		total := 0
		for _, id := range ids {
			total += huerfanos[id]
		}
		log.Printf("catalogo: %d compradores (%d proyectos fuera del catalogo) sin catalogo asociado, se conservan para el kNN: %v",
			total, len(ids), ids)
	}
	for _, afiliado := range afiliados {
		if afiliado.Cedula == "" {
			return nil, fmt.Errorf("catalogo: afiliado has empty cedula")
		}
		if _, exists := snapshot.afiliados[afiliado.Cedula]; exists {
			return nil, fmt.Errorf("catalogo: duplicate cedula %q", afiliado.Cedula)
		}
		snapshot.afiliados[afiliado.Cedula] = afiliado
	}

	entries, err := fs.ReadDir(data, "brochures")
	if err != nil {
		return nil, fmt.Errorf("catalogo: brochures: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || path.Ext(entry.Name()) != ".md" {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".md")
		if _, exists := snapshot.brochures[id]; exists {
			return nil, fmt.Errorf("catalogo: duplicate brochure %q", id)
		}
		content, err := fs.ReadFile(data, path.Join("brochures", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("catalogo: brochure %q: %w", id, err)
		}
		identity, err := brochureIdentity(string(content))
		if err != nil {
			return nil, fmt.Errorf("catalogo: brochure %q: %w", id, err)
		}
		if identity != id {
			return nil, fmt.Errorf("catalogo: brochure %q identity is %q", id, identity)
		}
		if _, exists := snapshot.proyectos[id]; !exists {
			return nil, fmt.Errorf("catalogo: brochure %q has no project", id)
		}
		snapshot.brochures[id] = string(content)
	}
	if len(snapshot.brochures) != len(snapshot.proyectos) {
		return nil, fmt.Errorf("catalogo: every project must have a brochure")
	}
	return &CatalogoRepository{snapshot: snapshot}, nil
}

func readJSON(data fs.FS, name string, target any) error {
	content, err := fs.ReadFile(data, name)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, target); err != nil {
		return err
	}
	return nil
}

func brochureIdentity(content string) (string, error) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "proyecto_id:") {
			id := strings.TrimSpace(strings.TrimPrefix(line, "proyecto_id:"))
			if id != "" {
				return id, nil
			}
		}
	}
	return "", fmt.Errorf("missing proyecto_id front matter")
}

func (r *CatalogoRepository) Proyectos(context.Context) (map[string]domain.Proyecto, error) {
	result := make(map[string]domain.Proyecto, len(r.snapshot.proyectos))
	for id, proyecto := range r.snapshot.proyectos {
		result[id] = proyecto
	}
	return result, nil
}

func (r *CatalogoRepository) Compradores(context.Context) ([]domain.Comprador, error) {
	return append([]domain.Comprador(nil), r.snapshot.compradores...), nil
}

func (r *CatalogoRepository) AfiliadoPorCedula(_ context.Context, cedula string) (*usecase.Afiliado, error) {
	afiliado, ok := r.snapshot.afiliados[cedula]
	if !ok {
		return nil, &usecase.NotFoundError{Resource: "afiliado", ID: cedula}
	}
	return &afiliado, nil
}

func (r *CatalogoRepository) BrochureMarkdown(_ context.Context, proyectoID string) (string, error) {
	brochure, ok := r.snapshot.brochures[proyectoID]
	if !ok {
		return "", &usecase.NotFoundError{Resource: "brochure", ID: proyectoID}
	}
	return brochure, nil
}
