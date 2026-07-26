package usecase

import (
	"context"
	"math"
	"reflect"
	"testing"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

func TestCalificarLeadPriorityCapsRatioAndUsesPartialConfidence(t *testing.T) {
	lead := calLead("priority-edge", false, domain.NivelAlta, 200)
	income := lead.Perfil["ingreso_hogar"]
	income.Fuente = domain.FuenteCampoDeclarado
	lead.Perfil["ingreso_hogar"] = income

	out, _, _, err := runCal(t, context.Background(), lead, calCatalog(100), nil)
	require(t, err == nil, "qualification failed: %v", err)
	require(t, out.Ruta == domain.RutaAsesor, "route = %q", out.Ruta)
	require(t, out.Capacidad.Ratio > 1.2, "ratio = %v, want strictly above 1.2", out.Capacidad.Ratio)
	require(t, out.Capacidad.Confianza < 1, "confidence = %v, want below 1", out.Capacidad.Confianza)
	wantPriority := 1.2 * out.Capacidad.Confianza
	require(t, out.Prioridad == wantPriority && out.Prioridad < 1.2,
		"priority = %v, want capped ratio times confidence = %v", out.Prioridad, wantPriority)
}

func TestCalificarLeadNonCanonicalCatalogKeyOmitsKNNZone(t *testing.T) {
	lead := calLead("zone-edge", false, domain.NivelAlta, 100)
	project := domain.Proyecto{ProyectoID: "p1", Zona: "NORTE", PrecioDesde: 100}
	buyer := domain.Comprador{ID: 1, ProyectoID: "p1", Categoria: "A", RangoEdad: "20-35", PersonasACargo: 1}

	withoutCanonicalZone := construirDecision(lead, map[string]domain.Proyecto{"display-name": project}, []domain.Comprador{buyer})
	withCanonicalZone := construirDecision(lead, map[string]domain.Proyecto{"p1": project}, []domain.Comprador{buyer})
	require(t, len(withoutCanonicalZone.vecinos) == 1 && len(withCanonicalZone.vecinos) == 1, "neighbor counts differ")

	withoutZone := withoutCanonicalZone.vecinos[0].Distancia
	withZone := withCanonicalZone.vecinos[0].Distancia
	wantWithoutZone := (0.15 * (7.5 / 32.5)) / 0.8
	wantWithZone := 0.15 * (7.5 / 32.5)
	require(t, math.Abs(withoutZone-wantWithoutZone) < 1e-12,
		"non-canonical zone distance = %.15f, want %.15f", withoutZone, wantWithoutZone)
	require(t, math.Abs(withZone-wantWithZone) < 1e-12 && withoutZone > withZone,
		"canonical/non-canonical distances = %.15f/%.15f", withZone, withoutZone)
}

func TestCalificarLeadRutaDecididaPayloadContainsCopiedDecisionData(t *testing.T) {
	catalog := calCatalog(100)
	catalog.buyers = []domain.Comprador{{ID: 1, ProyectoID: "p1", Categoria: "A", RangoEdad: "20-35", PersonasACargo: 1}}

	out, bus, _, err := runCal(t, context.Background(), calLead("event-edge", false, domain.NivelAlta, 120), catalog, nil)
	require(t, err == nil && len(bus.events) == 1, "qualification/events = %v/%d", err, len(bus.events))
	payload := bus.events[0].Payload
	require(t, payload["ruta"] == out.Ruta && payload["semaforo"] == out.Semaforo &&
		payload["consume_cupo_10"] == out.ConsumeCupo10 && payload["prioridad"] == out.Prioridad,
		"payload = %v", payload)

	recommendations, ok := payload["recomendaciones"].([]domain.Recomendacion)
	require(t, ok && reflect.DeepEqual(recommendations, out.Recomendaciones),
		"payload recommendations = %#v, output = %#v", payload["recomendaciones"], out.Recomendaciones)
	recommendations[0].Nombre = "mutated payload"
	require(t, out.Recomendaciones[0].Nombre != "mutated payload", "payload recommendations alias output")
}
