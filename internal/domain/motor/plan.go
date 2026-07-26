package motor

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

const maxCalendarMilestones = 4

// DisenarHitos creates the ordered, calendar-anchored milestones for a plan.
func DisenarHitos(brecha int64, convertible bool, desde time.Time, calendario []domain.EventoCalendario) []domain.Hito {
	start := desde.UTC()
	reassessment := dateOnly(start.AddDate(1, 0, 0))
	result := make([]domain.Hito, 0, maxCalendarMilestones+2)
	lower := dateOnly(start)
	if convertible {
		affiliation := lower.AddDate(0, 0, 8)
		result = append(result, domain.Hito{Tipo: domain.TipoHitoAfiliacion, Fecha: affiliation.Format("2006-01-02"), Descripcion: "Completar afiliación para activar beneficios", Estado: domain.EstadoHitoPendiente})
		lower = affiliation
	}

	type scheduled struct {
		event domain.EventoCalendario
		date  time.Time
		order int
	}
	scheduledEvents := make([]scheduled, 0, len(calendario))
	for order, event := range calendario {
		if event.Tipo != string(domain.TipoHitoCesantias) && event.Tipo != string(domain.TipoHitoPrima) {
			continue
		}
		date, ok := recurringDate(event.Fecha, start.Year(), lower, convertible)
		for occurrence := 0; ok && occurrence < maxCalendarMilestones && date.Before(reassessment); occurrence++ {
			scheduledEvents = append(scheduledEvents, scheduled{event: event, date: date, order: order})
			date = date.AddDate(1, 0, 0)
		}
	}
	sort.SliceStable(scheduledEvents, func(i, j int) bool {
		if scheduledEvents[i].date.Equal(scheduledEvents[j].date) {
			return scheduledEvents[i].order < scheduledEvents[j].order
		}
		return scheduledEvents[i].date.Before(scheduledEvents[j].date)
	})

	if brecha < 0 {
		brecha = 0
	}
	remaining := brecha
	for i := 0; i < len(scheduledEvents) && i < maxCalendarMilestones && remaining > 0; i++ {
		item := scheduledEvents[i]
		amount := remaining / 2
		if amount < 500_000 {
			amount = remaining
		}
		remaining -= amount
		monto := amount
		result = append(result, domain.Hito{Tipo: domain.TipoHito(item.event.Tipo), Fecha: item.date.Format("2006-01-02"), Monto: &monto, Descripcion: "Aporte asociado a " + strings.ToLower(item.event.Tipo), Estado: domain.EstadoHitoPendiente})
	}

	result = append(result, domain.Hito{Tipo: domain.TipoHitoReevaluacion, Fecha: reassessment.Format("2006-01-02"), Descripcion: "Revisar nuevamente la capacidad de compra", Estado: domain.EstadoHitoPendiente})
	return result
}

func recurringDate(value string, year int, lower time.Time, strictlyAfter bool) (time.Time, bool) {
	parts := strings.Split(strings.TrimPrefix(strings.TrimSpace(value), "--"), "-")
	if len(parts) != 2 {
		return time.Time{}, false
	}
	month, errMonth := strconv.Atoi(parts[0])
	day, errDay := strconv.Atoi(parts[1])
	if errMonth != nil || errDay != nil || month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}, false
	}
	candidate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if candidate.Month() != time.Month(month) || candidate.Day() != day {
		return time.Time{}, false
	}
	boundary := dateOnly(lower)
	if (strictlyAfter && !candidate.After(boundary)) || (!strictlyAfter && candidate.Before(boundary)) {
		candidate = candidate.AddDate(1, 0, 0)
	}
	return candidate, true
}

func dateOnly(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}
