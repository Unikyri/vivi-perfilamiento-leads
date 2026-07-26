package motor

import (
	"math"
	"sort"
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

// plegarAcento reduce las vocales acentuadas y la eñe del español a ASCII.
// Se resuelve con un switch en vez de golang.org/x/text para que el paquete
// domain siga sin dependencias externas (NFR-M-01).
func plegarAcento(r rune) rune {
	switch r {
	case 'á', 'à', 'ä', 'â':
		return 'a'
	case 'é', 'è', 'ë', 'ê':
		return 'e'
	case 'í', 'ì', 'ï', 'î':
		return 'i'
	case 'ó', 'ò', 'ö', 'ô':
		return 'o'
	case 'ú', 'ù', 'ü', 'û':
		return 'u'
	case 'ñ':
		return 'n'
	}
	return r
}

// palabrasVaciasZona son términos de mercadeo que aparecen en el nombre
// comercial de la zona pero no identifican un lugar. Sin esto, "Ciudadela
// Maiporé - Soacha" y "Ciudadela Calle 80 - Bogotá" coincidirían por la
// palabra "ciudadela", que es exactamente lo que no queremos.
var palabrasVaciasZona = map[string]bool{
	"ciudadela": true, "conjunto": true, "agrupacion": true,
	"residencial": true, "vivienda": true, "urbanizacion": true,
	"de": true, "del": true, "la": true, "el": true,
	"los": true, "las": true, "y": true,
}

// zonaTokens reduce una zona a su conjunto de tokens identificadores.
func zonaTokens(zona string) map[string]bool {
	plegado := strings.Map(plegarAcento, strings.ToLower(normalizeOptionalText(zona)))
	plegado = strings.NewReplacer("-", " ", ",", " ", ".", " ", "/", " ").Replace(plegado)

	tokens := make(map[string]bool)
	for _, token := range strings.Fields(plegado) {
		if palabrasVaciasZona[token] {
			continue
		}
		tokens[token] = true
	}
	return tokens
}

// zonasCoinciden compara zonas por intersección de tokens identificadores,
// no por igualdad exacta: el lead escribe "Soacha" y el catálogo dice
// "Ciudadela Maiporé - Soacha".
//
// El recorrido del mapa es no determinista, pero el resultado no depende del
// orden: es una intersección booleana. NO introducir aquí lógica que dependa
// del orden (primer match ganador con retorno de valor, acumulación en slice).
func zonasCoinciden(izquierda, derecha string) bool {
	tokensIzq := zonaTokens(izquierda)
	if len(tokensIzq) == 0 {
		return false
	}
	tokensDer := zonaTokens(derecha)
	if len(tokensDer) == 0 {
		return false
	}
	for token := range tokensIzq {
		if tokensDer[token] {
			return true
		}
	}
	return false
}

const (
	gowerCategoryWeight    = 0.35
	gowerZoneWeight        = 0.20
	gowerAgeWeight         = 0.15
	gowerAffiliationWeight = 0.15
	gowerDependentsWeight  = 0.15
	gowerAgeRange          = 32.5
	gowerDependentsRange   = 10.0
)

type projectAffiliateStats struct {
	affiliateCount int
	categoryMode   string
	hasCategory    bool
	ageMedian      float64
	hasAge         bool
}

type projectAffiliateStatsAccumulator struct {
	affiliateCount int
	categoryCounts map[string]int
	ages           []float64
}

// buildProjectAffiliateStats creates a value-only index from affiliated
// historical buyers. Statistics are keyed by exact project ID and are not
// mutated after this function returns.
func buildProjectAffiliateStats(buyers []domain.Comprador) map[string]projectAffiliateStats {
	accumulators := make(map[string]*projectAffiliateStatsAccumulator)
	for _, buyer := range buyers {
		if !buyer.Afiliado || buyer.ProyectoID == "" {
			continue
		}

		stats := accumulators[buyer.ProyectoID]
		if stats == nil {
			stats = &projectAffiliateStatsAccumulator{categoryCounts: make(map[string]int)}
			accumulators[buyer.ProyectoID] = stats
		}
		stats.affiliateCount++
		if category := normalizeCategory(buyer.Categoria); category != "" {
			stats.categoryCounts[category]++
		}
		if age, ok := ageRepresentative(buyer.RangoEdad); ok {
			stats.ages = append(stats.ages, age)
		}
	}

	result := make(map[string]projectAffiliateStats, len(accumulators))
	for projectID, accumulated := range accumulators {
		stats := projectAffiliateStats{affiliateCount: accumulated.affiliateCount}
		if accumulated.affiliateCount >= 5 {
			stats.categoryMode, stats.hasCategory = categoryMode(accumulated.categoryCounts)
			stats.ageMedian, stats.hasAge = numericMedian(accumulated.ages)
		}
		result[projectID] = stats
	}
	return result
}

func categoryMode(counts map[string]int) (string, bool) {
	best := ""
	bestCount := 0
	for _, category := range []string{"A", "B", "C"} {
		if counts[category] > bestCount {
			best = category
			bestCount = counts[category]
		}
	}
	return best, bestCount > 0
}

func numericMedian(values []float64) (float64, bool) {
	if len(values) == 0 {
		return 0, false
	}
	ordered := append([]float64(nil), values...)
	sort.Float64s(ordered)
	middle := len(ordered) / 2
	if len(ordered)%2 == 1 {
		return ordered[middle], true
	}
	return (ordered[middle-1] + ordered[middle]) / 2, true
}

// projectBuyerWithStats applies category and age fallback independently. Only
// non-affiliates are imputed; affiliates are the provenance source.
func projectBuyerWithStats(buyer domain.Comprador, catalogZones map[string]string, projectStats map[string]projectAffiliateStats) featureVector {
	features := projectBuyer(buyer, catalogZones)
	if buyer.Afiliado {
		return features
	}

	stats, ok := projectStats[buyer.ProyectoID]
	if !ok || stats.affiliateCount < 5 {
		return features
	}
	if !features.tieneCategoria && stats.hasCategory {
		features.categoria = stats.categoryMode
		features.tieneCategoria = true
	}
	if !features.tieneEdad && stats.hasAge {
		features.edad = stats.ageMedian
		features.tieneEdad = true
	}
	return features
}

// gowerDistance returns the weighted distance over mutually present
// dimensions. Missing dimensions do not contribute weight, so the remaining
// dimensions are renormalized by the participating weight total.
func gowerDistance(left, right featureVector) float64 {
	weightedDistance := 0.0
	participatingWeight := 0.0
	add := func(weight, distance float64) {
		weightedDistance += weight * distance
		participatingWeight += weight
	}

	if left.tieneCategoria && right.tieneCategoria {
		if left.categoria == right.categoria {
			add(gowerCategoryWeight, 0)
		} else {
			add(gowerCategoryWeight, 1)
		}
	}
	if left.tieneZona && right.tieneZona {
		if zonasCoinciden(left.zona, right.zona) {
			add(gowerZoneWeight, 0)
		} else {
			add(gowerZoneWeight, 1)
		}
	}
	if left.tieneEdad && right.tieneEdad {
		add(gowerAgeWeight, cappedNumericDistance(left.edad, right.edad, gowerAgeRange))
	}
	if left.tieneAfiliado && right.tieneAfiliado {
		if left.afiliado == right.afiliado {
			add(gowerAffiliationWeight, 0)
		} else {
			add(gowerAffiliationWeight, 1)
		}
	}
	if left.tienePersonasACargo && right.tienePersonasACargo {
		add(gowerDependentsWeight, cappedNumericDistance(
			float64(left.personasACargo), float64(right.personasACargo), gowerDependentsRange,
		))
	}

	if participatingWeight == 0 {
		return 0
	}
	return weightedDistance / participatingWeight
}

func cappedNumericDistance(left, right, valueRange float64) float64 {
	if valueRange <= 0 {
		return 0
	}
	distance := math.Abs(left-right) / valueRange
	if distance > 1 {
		return 1
	}
	return distance
}

// GemeloKNN returns the closest buyers as value-only neighbors. Selection is
// deterministic: distance ascending, then buyer ID ascending. The input
// records and catalog are never modified.
func GemeloKNN(in EntradaGemelo) []Vecino {
	if in.K <= 0 || len(in.Compradores) == 0 {
		return make([]Vecino, 0)
	}

	lead := projectLead(in)
	projectStats := buildProjectAffiliateStats(in.Compradores)
	neighbors := make([]Vecino, 0, len(in.Compradores))
	for _, buyer := range in.Compradores {
		features := projectBuyerWithStats(buyer, in.ZonasCatalogo, projectStats)
		neighbors = append(neighbors, Vecino{
			ID:         buyer.ID,
			ProyectoID: buyer.ProyectoID,
			Desistio:   buyer.Desistio,
			Distancia:  gowerDistance(lead, features),
		})
	}

	sort.SliceStable(neighbors, func(i, j int) bool {
		if neighbors[i].Distancia != neighbors[j].Distancia {
			return neighbors[i].Distancia < neighbors[j].Distancia
		}
		return neighbors[i].ID < neighbors[j].ID
	})
	if in.K < len(neighbors) {
		neighbors = neighbors[:in.K]
	}
	return neighbors
}
