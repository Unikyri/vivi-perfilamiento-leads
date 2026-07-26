package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

const demoSeedTimeout = 3 * time.Second

var ErrDemoSeedDeshabilitado = errors.New("demo seed disabled")

// CargarSeed is intentionally inert unless the operator explicitly enables it.
type CargarSeed struct {
	Repository DemoSeedRepository
	Habilitado bool
}

func (uc *CargarSeed) Ejecutar(ctx context.Context) error {
	if uc == nil || !uc.Habilitado {
		return nil
	}
	if uc.Repository == nil {
		return fmt.Errorf("%w: repositorio de seed requerido", ErrValidacion)
	}
	seedCtx, cancel := context.WithTimeout(ctx, demoSeedTimeout)
	defer cancel()
	if err := uc.Repository.Sembrar(seedCtx, SemillasDemo()); err != nil {
		return fmt.Errorf("cargar seed demo: %w", err)
	}
	return nil
}

// FechaDemoAprobada is the deterministic date used by the controlled demo.
func FechaDemoAprobada() time.Time {
	return time.Date(2026, 7, 26, 0, 0, 0, 0, time.UTC)
}

// SemillasDemo returns fresh values so callers cannot mutate the canonical set.
func SemillasDemo() []domain.Lead {
	now := FechaDemoAprobada()
	return []domain.Lead{
		seedLead("ana", "Ana Rodríguez", "+57 300 123 4567", "1032456789", true, 0.90, now),
		seedLead("carlos", "Carlos Martínez", "+57 311 987 6543", "1000000000", false, 0.80, now),
		seedLead("luisa", "Luisa Gómez", "+57 300 000 0000", "1098765432", true, 0.70, now),
	}
}

func seedLead(id, name, phone, document string, affiliate bool, priority float64, now time.Time) domain.Lead {
	profile := domain.Perfil{
		"ingreso_hogar":    {Valor: int64(3500000), Fuente: domain.FuenteCampoDeclarado, Confianza: 1, ActualizadoEn: now},
		"recursos_propios": {Valor: int64(10000000), Fuente: domain.FuenteCampoDeclarado, Confianza: 1, ActualizadoEn: now},
		"personas_hogar":   {Valor: int64(2), Fuente: domain.FuenteCampoDeclarado, Confianza: 1, ActualizadoEn: now},
		"tiene_vivienda":   {Valor: false, Fuente: domain.FuenteCampoDeclarado, Confianza: 1, ActualizadoEn: now},
		"recibio_subsidio": {Valor: false, Fuente: domain.FuenteCampoDeclarado, Confianza: 1, ActualizadoEn: now},
		"categoria":        {Valor: "A", Fuente: domain.FuenteCampoDeclarado, Confianza: 1, ActualizadoEn: now},
	}
	return domain.Lead{
		LeadID: id, Nombre: name, Telefono: phone, Cedula: document, Fuente: "DEMO",
		Estado: domain.EstadoLeadPerfilando, Ruta: domain.RutaAsesor, Afiliado: affiliate,
		Prioridad: priority, ConsumeCupo10: !affiliate, Perfil: profile,
		Capacidad: &domain.Capacidad{PresupuestoMax: 180000000, RecursosPropios: 10000000, Confianza: 1},
		Version:   1, CreadoEn: now, ActualizadoEn: now,
	}
}
