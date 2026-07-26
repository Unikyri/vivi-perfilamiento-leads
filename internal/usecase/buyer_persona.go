package usecase

import (
	"context"
	"sort"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

var buyerCategories = []string{"A", "B", "C", "SIN_DATO"}
var buyerAgeBands = []string{"20-35", "36-45", "46-55", "55+", "SIN_DATO"}

type AfiliacionPersona struct {
	Afiliados   float64 `json:"afiliados"`
	NoAfiliados float64 `json:"no_afiliados"`
}

type BuyerPersonaResumen struct {
	ProyectoID        string             `json:"proyecto_id"`
	Nombre            string             `json:"nombre"`
	Muestras          int                `json:"muestras"`
	Afiliacion        AfiliacionPersona  `json:"afiliacion"`
	Categoria         map[string]float64 `json:"categoria"`
	RangoEdad         map[string]float64 `json:"rango_edad"`
	TasaDesistimiento float64            `json:"tasa_desistimiento"`
	ActualizadoEn     string             `json:"actualizado_en"`
}

type BuyerPersonaCatalogo struct {
	Proyectos []BuyerPersonaResumen `json:"proyectos"`
}

type BuyerPersona struct {
	Catalogo CatalogoRepository
}

func (b *BuyerPersona) Proyecto(ctx context.Context, proyectoID string) (BuyerPersonaResumen, error) {
	if b == nil || b.Catalogo == nil || proyectoID == "" {
		return BuyerPersonaResumen{}, ErrValidacion
	}
	proyectos, err := b.Catalogo.Proyectos(ctx)
	if err != nil {
		return BuyerPersonaResumen{}, err
	}
	proyecto, ok := proyectos[proyectoID]
	if !ok {
		return BuyerPersonaResumen{}, &NotFoundError{Resource: "proyecto", ID: proyectoID}
	}
	compradores, err := b.Catalogo.Compradores(ctx)
	if err != nil {
		return BuyerPersonaResumen{}, err
	}
	return resumirBuyerPersona(proyecto, compradores), nil
}

func (b *BuyerPersona) CatalogoCompleto(ctx context.Context) (BuyerPersonaCatalogo, error) {
	if b == nil || b.Catalogo == nil {
		return BuyerPersonaCatalogo{}, ErrValidacion
	}
	proyectos, err := b.Catalogo.Proyectos(ctx)
	if err != nil {
		return BuyerPersonaCatalogo{}, err
	}
	compradores, err := b.Catalogo.Compradores(ctx)
	if err != nil {
		return BuyerPersonaCatalogo{}, err
	}
	ids := make([]string, 0, len(proyectos))
	for id := range proyectos {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	result := BuyerPersonaCatalogo{Proyectos: make([]BuyerPersonaResumen, 0, len(ids))}
	for _, id := range ids {
		result.Proyectos = append(result.Proyectos, resumirBuyerPersona(proyectos[id], compradores))
	}
	return result, nil
}

func resumirBuyerPersona(proyecto domain.Proyecto, compradores []domain.Comprador) BuyerPersonaResumen {
	categoria := make(map[string]float64, len(buyerCategories))
	for _, key := range buyerCategories {
		categoria[key] = 0
	}
	edad := make(map[string]float64, len(buyerAgeBands))
	for _, key := range buyerAgeBands {
		edad[key] = 0
	}
	var total, afiliados, desistidos int
	var latest time.Time
	for _, comprador := range compradores {
		if comprador.ProyectoID != proyecto.ProyectoID {
			continue
		}
		total++
		if comprador.Afiliado {
			afiliados++
		}
		if comprador.Desistio {
			desistidos++
		}
		if _, ok := categoria[comprador.Categoria]; ok {
			categoria[comprador.Categoria]++
		}
		if _, ok := edad[comprador.RangoEdad]; ok {
			edad[comprador.RangoEdad]++
		}
		if date, ok := fechaFuente(comprador.FechaOpcion); ok && date.After(latest) {
			latest = date
		}
	}
	if total > 0 {
		for key := range categoria {
			categoria[key] /= float64(total)
		}
		for key := range edad {
			edad[key] /= float64(total)
		}
	}
	denom := float64(total)
	if total == 0 {
		denom = 1
	}
	return BuyerPersonaResumen{
		ProyectoID: proyecto.ProyectoID,
		Nombre:     proyecto.Nombre,
		Muestras:   total,
		Afiliacion: AfiliacionPersona{
			Afiliados:   float64(afiliados) / denom,
			NoAfiliados: float64(total-afiliados) / denom,
		},
		Categoria:         categoria,
		RangoEdad:         edad,
		TasaDesistimiento: float64(desistidos) / denom,
		ActualizadoEn:     latest.UTC().Format(time.RFC3339),
	}
}

func fechaFuente(raw string) (time.Time, bool) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02"} {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return parsed.UTC(), true
		}
	}
	return time.Time{}, false
}
