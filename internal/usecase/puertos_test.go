package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
)

type leadRepositoryShape struct{}

func (leadRepositoryShape) Crear(context.Context, *domain.Lead) error { return nil }
func (leadRepositoryShape) PorID(context.Context, string) (*domain.Lead, error) {
	return nil, nil
}
func (leadRepositoryShape) Guardar(context.Context, *domain.Lead) error { return nil }
func (leadRepositoryShape) Listar(context.Context, FiltroLeads) ([]*domain.Lead, error) {
	return nil, nil
}
func (leadRepositoryShape) AgregarMensaje(context.Context, *domain.Mensaje) error { return nil }
func (leadRepositoryShape) Conversacion(context.Context, string) ([]domain.Mensaje, error) {
	return nil, nil
}

var _ LeadRepository = leadRepositoryShape{}

func TestPuertos_MethodShapes(t *testing.T) {
	var (
		_ PlanRepository     = planRepositoryShape{}
		_ FichaRepository    = fichaRepositoryShape{}
		_ CatalogoRepository = catalogoRepositoryShape{}
		_ LLMProvider        = llmProviderShape{}
		_ MensajeriaGateway  = mensajeriaGatewayShape{}
		_ Reloj              = relojShape{}
		_ BusEventos         = busEventosShape{}
		_ GeneradorID        = generadorIDShape{}
	)
}

type planRepositoryShape struct{}

func (planRepositoryShape) Crear(context.Context, *domain.PlanNutricion) error { return nil }
func (planRepositoryShape) PorLead(context.Context, string) (*domain.PlanNutricion, error) {
	return nil, nil
}
func (planRepositoryShape) Guardar(context.Context, *domain.PlanNutricion) error { return nil }
func (planRepositoryShape) HitosVencidos(context.Context, time.Time) ([]HitoConPlan, error) {
	return nil, nil
}
func (planRepositoryShape) MarcarHito(context.Context, string, domain.EstadoHito) error { return nil }

type fichaRepositoryShape struct{}

func (fichaRepositoryShape) Guardar(context.Context, *domain.Ficha) error           { return nil }
func (fichaRepositoryShape) PorLead(context.Context, string) (*domain.Ficha, error) { return nil, nil }

type catalogoRepositoryShape struct{}

func (catalogoRepositoryShape) Proyectos(context.Context) (map[string]domain.Proyecto, error) {
	return nil, nil
}
func (catalogoRepositoryShape) Compradores(context.Context) ([]domain.Comprador, error) {
	return nil, nil
}
func (catalogoRepositoryShape) AfiliadoPorCedula(context.Context, string) (*Afiliado, error) {
	return nil, nil
}
func (catalogoRepositoryShape) BrochureMarkdown(context.Context, string) (string, error) {
	return "", nil
}

type llmProviderShape struct{}

func (llmProviderShape) GenerarTurno(context.Context, EntradaTurno) (SalidaTurno, error) {
	return SalidaTurno{}, nil
}
func (llmProviderShape) ProcesarAudio(context.Context, Audio, EntradaTurno) (SalidaTurno, error) {
	return SalidaTurno{}, nil
}
func (llmProviderShape) Nombre() string { return "test" }

type mensajeriaGatewayShape struct{}

func (mensajeriaGatewayShape) Enviar(context.Context, *domain.Mensaje) error { return nil }

type relojShape struct{}

func (relojShape) Ahora() time.Time         { return time.Time{} }
func (relojShape) FechaSimulada() time.Time { return time.Time{} }
func (relojShape) Avanzar(time.Time)        {}

type busEventosShape struct{}

func (busEventosShape) Publicar(context.Context, Evento)                {}
func (busEventosShape) Suscribir(string, func(context.Context, Evento)) {}

type generadorIDShape struct{}

func (generadorIDShape) Nuevo() string { return "id" }

func TestDTO_JSONContract(t *testing.T) {
	falseValue := false
	cases := []struct {
		name  string
		value any
		keys  []string
	}{
		{
			name:  "affiliate",
			value: Afiliado{Cedula: "1", Nombre: "Ana", Categoria: "A", IngresoMensual: 2600000},
			keys:  []string{"cedula", "nombre", "categoria", "ingreso_mensual"},
		},
		{
			name:  "turn input",
			value: EntradaTurno{LeadID: "lead-1", Nombre: "Ana", MensajeUsuario: "hola", NumerosDelMotor: map[string]int64{"presupuesto_max": 100}, EsAfiliado: true},
			keys:  []string{"lead_id", "nombre", "mensaje_usuario", "numeros_del_motor", "es_afiliado"},
		},
		{
			name: "llm output",
			value: SalidaTurno{
				CamposExtraidos: []CampoExtraido{{Campo: "recursos_propios", Valor: 8000000, Fuente: domain.FuenteCampoDeclarado, Confianza: .85, RequiereConfirmacion: falseValue}},
				Intencion:       domain.Intencion{Nivel: domain.NivelAlta}, Respuesta: "continua", Accion: "CONTINUAR",
			},
			keys: []string{"campos_extraidos", "intencion", "respuesta", "accion"},
		},
		{
			name:  "event",
			value: Evento{Tipo: EvLeadNuevo, LeadID: "lead-1", Payload: map[string]any{"nombre": "Ana"}},
			keys:  []string{"tipo", "lead_id", "payload"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			var object map[string]any
			if err := json.Unmarshal(encoded, &object); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			for _, key := range tc.keys {
				if _, ok := object[key]; !ok {
					t.Errorf("JSON missing %q: %s", key, encoded)
				}
			}
		})
	}
}

func TestErrores_MatchWrappedSentinels(t *testing.T) {
	wrappedNotFound := &NotFoundError{Resource: "lead", ID: "missing"}
	if !errors.Is(wrappedNotFound, ErrNoEncontrado) {
		t.Fatalf("errors.Is(%v, ErrNoEncontrado) = false", wrappedNotFound)
	}
	wrappedLock := errors.Join(errors.New("guardar"), ErrOptimisticLock)
	if !errors.Is(wrappedLock, ErrOptimisticLock) {
		t.Fatalf("errors.Is(%v, ErrOptimisticLock) = false", wrappedLock)
	}
}

func TestEventoConstants_ContractNames(t *testing.T) {
	cases := map[string]string{
		EvLeadNuevo: "LeadNuevo", EvMensajeEntrante: "MensajeEntrante", EvPerfilCompleto: "PerfilCompleto",
		EvRutaDecidida: "RutaDecidida", EvHitoVencido: "HitoVencido", EvTickReloj: "TickReloj", EvResultadoReal: "ResultadoReal",
	}
	if len(cases) != 7 {
		t.Fatalf("event constant count = %d, want 7", len(cases))
	}
	for got, want := range cases {
		if got != want {
			t.Errorf("event constant = %q, want %q", got, want)
		}
	}
}
