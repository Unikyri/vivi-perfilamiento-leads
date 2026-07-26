package usecase

import (
	"context"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"strings"
	"testing"
	"time"
)

func TestSaludarLeadDeterministicGreeting(t *testing.T) {
	for _, affiliated := range []bool{false, true} {
		repo := NuevoLeadRepoFake()
		lead := &domain.Lead{LeadID: "lead-1", Nombre: "Ana", Afiliado: affiliated}
		if err := repo.Crear(context.Background(), lead); err != nil {
			t.Fatal(err)
		}
		uc := &SaludarLead{Leads: repo, IDs: NuevoIDFake("msg"), Reloj: NuevoRelojFake(time.Unix(10, 0))}
		if err := uc.Ejecutar(context.Background(), Evento{Tipo: EvLeadNuevo, LeadID: lead.LeadID}); err != nil {
			t.Fatal(err)
		}
		messages, err := repo.Conversacion(context.Background(), lead.LeadID)
		if err != nil {
			t.Fatal(err)
		}
		if len(messages) != 1 || messages[0].Autor != domain.AutorMensajeVivi || messages[0].TipoContenido != domain.TipoContenidoTexto || !strings.Contains(messages[0].Texto, "Ana") {
			t.Fatalf("unexpected greeting: %+v", messages)
		}
		if affiliated != strings.Contains(messages[0].Texto, "subsidio") {
			t.Fatalf("affiliate greeting mismatch: %q", messages[0].Texto)
		}
	}
}
