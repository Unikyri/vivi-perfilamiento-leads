package agentes

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

var errSkipped = errors.New("event skipped")

// Coordinadora installs the fixed, provider-free event handoff table.
type Coordinadora struct {
	bus  usecase.BusEventos
	deps Dependencias
	once sync.Once
}

// Nueva constructs a coordinator. A nil bus is safe and simply disables wiring.
func Nueva(bus usecase.BusEventos, deps Dependencias) *Coordinadora {
	return &Coordinadora{bus: bus, deps: deps}
}

// Registrar installs each available consumer once, in the canonical event order.
func (c *Coordinadora) Registrar() {
	if c == nil || c.bus == nil {
		return
	}
	c.once.Do(func() {
		if c.deps.LeadNuevo != nil {
			c.subscribe(usecase.EvLeadNuevo, "lead_nuevo", c.deps.LeadNuevo)
		}
		if c.deps.Calificador != nil {
			c.subscribe(usecase.EvPerfilCompleto, "calificar_lead", c.perfilCompleto)
		}
		if c.deps.Documentadora != nil {
			c.subscribe(usecase.EvRutaDecidida, "generar_ficha", c.rutaDecidida)
		}
		if c.deps.Nutricionista != nil {
			c.subscribe(usecase.EvTickReloj, "ejecutar_hitos", c.tickReloj)
		}
	})
}

func (c *Coordinadora) subscribe(tipo, nombre string, handler ManejadorEvento) {
	c.bus.Suscribir(tipo, func(ctx context.Context, event usecase.Evento) {
		err := handler(ctx, event)
		outcome := "ok"
		if err != nil {
			outcome = "error"
			if errors.Is(err, errSkipped) || errors.Is(err, usecase.ErrLeadNoCalificable) {
				outcome = "skipped"
			}
		}
		c.log(event, nombre, outcome)
	})
}

func (c *Coordinadora) log(event usecase.Evento, handler, outcome string) {
	if c.deps.Logger == nil {
		return
	}
	c.deps.Logger.Info("coordinator dispatch",
		slog.String("tipo", event.Tipo),
		slog.String("lead_id", event.LeadID),
		slog.String("handler", handler),
		slog.String("outcome", outcome),
	)
}

func (c *Coordinadora) perfilCompleto(ctx context.Context, event usecase.Evento) error {
	if strings.TrimSpace(event.LeadID) == "" {
		return errSkipped
	}
	_, err := c.deps.Calificador.Ejecutar(ctx, usecase.EntradaCalificar{LeadID: event.LeadID})
	if errors.Is(err, usecase.ErrLeadNoCalificable) {
		return errSkipped
	}
	return err
}

func (c *Coordinadora) rutaDecidida(ctx context.Context, event usecase.Evento) error {
	if strings.TrimSpace(event.LeadID) == "" {
		return errSkipped
	}
	ruta, ok := routeFromPayload(event.Payload)
	if !ok || ruta != domain.RutaAsesor {
		return errSkipped
	}
	_, err := c.deps.Documentadora.Ejecutar(ctx, event.LeadID)
	return err
}

func (c *Coordinadora) tickReloj(ctx context.Context, event usecase.Evento) error {
	hasta, ok := timeFromPayload(event.Payload)
	if !ok {
		return errSkipped
	}
	count, err := c.deps.Nutricionista.Ejecutar(ctx, hasta)
	if result := usecase.ResultadoTickDeContexto(ctx); result != nil {
		result.HitosDisparados = count
		result.Err = err
	}
	return err
}

func routeFromPayload(payload map[string]any) (domain.Ruta, bool) {
	if payload == nil {
		return "", false
	}
	switch value := payload["ruta"].(type) {
	case domain.Ruta:
		return value, value != ""
	case string:
		route := domain.Ruta(strings.TrimSpace(value))
		return route, route != ""
	default:
		return "", false
	}
}

func timeFromPayload(payload map[string]any) (time.Time, bool) {
	if payload == nil {
		return time.Time{}, false
	}
	for _, key := range []string{"fecha_simulada", "hasta"} {
		switch value := payload[key].(type) {
		case time.Time:
			return value, !value.IsZero()
		case string:
			parsed, err := time.Parse(time.RFC3339, value)
			if err == nil {
				return parsed, true
			}
		}
	}
	return time.Time{}, false
}
