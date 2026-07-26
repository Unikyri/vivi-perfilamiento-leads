package motor

import (
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"reflect"
	"testing"
	"time"
)

func TestDisenarHitosDeterministicCalendarPlan(t *testing.T) {
	calendar := []domain.EventoCalendario{{Tipo: "PRIMA", Fecha: "--12-20"}, {Tipo: "CESANTIAS", Fecha: "--02-14"}, {Tipo: "PRIMA", Fecha: "--06-30"}}
	start := time.Date(2026, 1, 10, 18, 0, 0, 0, time.FixedZone("local", -5*60*60))
	before := append([]domain.EventoCalendario(nil), calendar...)
	first := DisenarHitos(1_000_001, false, start, calendar)
	second := DisenarHitos(1_000_001, false, start, calendar)
	if !reflect.DeepEqual(first, second) || !reflect.DeepEqual(calendar, before) {
		t.Fatalf("planner is not deterministic or mutated input: %#v", first)
	}
	if len(first) != 3 || first[0].Tipo != domain.TipoHitoCesantias || first[0].Fecha != "2026-02-14" || first[1].Fecha != "2026-06-30" || first[2].Tipo != domain.TipoHitoReevaluacion {
		t.Fatalf("calendar order = %#v", first)
	}
	if first[0].Monto == nil || *first[0].Monto != 500_000 || first[1].Monto == nil || *first[1].Monto != 500_001 {
		t.Fatalf("gap allocation = %#v", first)
	}
}
func TestDisenarHitosConversionFirstAndRollover(t *testing.T) {
	calendar := []domain.EventoCalendario{{Tipo: "PRIMA", Fecha: "--01-15"}, {Tipo: "CESANTIAS", Fecha: "--02-01"}, {Tipo: "PRIMA", Fecha: "--12-01"}}
	got := DisenarHitos(0, true, time.Date(2026, 12, 20, 23, 0, 0, 0, time.UTC), calendar)
	if len(got) != 2 || got[0].Tipo != domain.TipoHitoAfiliacion || got[0].Fecha != "2026-12-28" || got[1].Tipo != domain.TipoHitoReevaluacion || got[1].Fecha != "2027-12-20" {
		t.Fatalf("conversion/rollover = %#v", got)
	}
}
func TestDisenarHitosSmallOddAndNegativeGaps(t *testing.T) {
	calendar := []domain.EventoCalendario{{Tipo: "PRIMA", Fecha: "--03-01"}}
	for _, tc := range []struct {
		name      string
		gap, want int64
		hasAmount bool
	}{{"small", 7, 7, true}, {"odd", 1_000_001, 500_000, true}, {"negative", -1, 0, false}} {
		t.Run(tc.name, func(t *testing.T) {
			got := DisenarHitos(tc.gap, false, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), calendar)
			if (got[0].Monto != nil) != tc.hasAmount || (tc.hasAmount && *got[0].Monto != tc.want) {
				t.Fatalf("first amount = %#v, want %d", got[0].Monto, tc.want)
			}
		})
	}
}
func TestDisenarHitosKeepsReevaluationLastByDate(t *testing.T) {
	calendar := []domain.EventoCalendario{{Tipo: "PRIMA", Fecha: "--12-15"}, {Tipo: "CESANTIAS", Fecha: "--01-15"}}
	got := DisenarHitos(4_000_004, false, time.Date(2026, 12, 20, 23, 0, 0, 0, time.FixedZone("local", -5*60*60)), calendar)
	wantDates := []string{"2027-01-15", "2027-12-15", "2027-12-21"}
	if len(got) != len(wantDates) {
		t.Fatalf("milestones = %#v, want two monetary milestones before reevaluation", got)
	}
	for i, want := range wantDates {
		if got[i].Fecha != want || (i < 2 && got[i].Monto == nil) {
			t.Fatalf("milestone %d = %#v, want %s", i, got[i], want)
		}
	}
}
