package motor

import (
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

const smmlv2026 int64 = 1_750_905

// EntradaGemelo is the complete, contract-safe input for buyer-twin projection.
// Catalog zones are resolved only by exact ProyectoID; callers cannot provide
// independently indexed feature slices or buyer-owned zone values.
type EntradaGemelo struct {
	Perfil         domain.Perfil
	Afiliado       bool
	PersonasACargo *int
	Compradores    []domain.Comprador
	ZonasCatalogo  map[string]string
	K              int
}

// Vecino is the value-only result shape reserved for the selection phase.
type Vecino struct {
	ID         int
	ProyectoID string
	Desistio   bool
	Distancia  float64
}

// featureVector is the internal projection shape. Presence is tracked per
// dimension because zero dependents and false affiliation are valid values.
type featureVector struct {
	categoria           string
	tieneCategoria      bool
	zona                string
	tieneZona           bool
	edad                float64
	tieneEdad           bool
	afiliado            bool
	tieneAfiliado       bool
	personasACargo      int
	tienePersonasACargo bool
}

func projectLead(in EntradaGemelo) featureVector {
	features := featureVector{
		afiliado:      in.Afiliado,
		tieneAfiliado: true,
	}

	if income, ok := in.Perfil.Entero("ingreso_hogar"); ok && income >= 0 {
		features.categoria = incomeCategory(income)
		features.tieneCategoria = true
	}

	if age, ok := in.Perfil.Entero("edad"); ok && age > 0 {
		features.edad = float64(age)
		features.tieneEdad = true
	}

	if zone, ok := in.Perfil.Texto("zona_deseada"); ok {
		if normalized := normalizeOptionalText(zone); normalized != "" {
			features.zona = normalized
			features.tieneZona = true
		}
	}

	if in.PersonasACargo != nil && *in.PersonasACargo >= 0 {
		features.personasACargo = *in.PersonasACargo
		features.tienePersonasACargo = true
	}

	return features
}

func projectBuyer(buyer domain.Comprador, catalogZones map[string]string) featureVector {
	features := featureVector{
		afiliado:       buyer.Afiliado,
		tieneAfiliado:  true,
		personasACargo: buyer.PersonasACargo,
	}
	if buyer.PersonasACargo >= 0 {
		features.tienePersonasACargo = true
	}

	if category := normalizeCategory(buyer.Categoria); category != "" {
		features.categoria = category
		features.tieneCategoria = true
	}

	if age, ok := ageRepresentative(buyer.RangoEdad); ok {
		features.edad = age
		features.tieneEdad = true
	}

	// ProyectoID is deliberately not normalized or inferred. A zone is valid
	// only when the caller supplied an exact catalog entry for that ID.
	if buyer.ProyectoID != "" {
		if zone, ok := catalogZones[buyer.ProyectoID]; ok {
			if normalized := normalizeOptionalText(zone); normalized != "" {
				features.zona = normalized
				features.tieneZona = true
			}
		}
	}

	return features
}

func incomeCategory(income int64) string {
	switch {
	case income <= 2*smmlv2026:
		return "A"
	case income <= 4*smmlv2026:
		return "B"
	default:
		return "C"
	}
}

func ageRepresentative(label string) (float64, bool) {
	switch normalizeOptionalText(label) {
	case "20-35":
		return 27.5, true
	case "36-45":
		return 40.5, true
	case "46-55":
		return 50.5, true
	case "55+":
		return 60.0, true
	default:
		return 0, false
	}
}

func normalizeCategory(value string) string {
	switch strings.ToUpper(normalizeOptionalText(value)) {
	case "A":
		return "A"
	case "B":
		return "B"
	case "C":
		return "C"
	default:
		return ""
	}
}

func normalizeOptionalText(value string) string {
	normalized := strings.Join(strings.Fields(value), " ")
	if strings.EqualFold(normalized, "SIN_DATO") {
		return ""
	}
	return normalized
}
