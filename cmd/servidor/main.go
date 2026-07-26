package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/adapters/agentes"
	adapterhttp "github.com/Unikyri/vivi-perfilamiento-leads/internal/adapters/http"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/domain"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/bus"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/config"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/ids"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/llm"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/postgres"
	infrareloj "github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/reloj"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

// version la inyecta el build: -ldflags "-X main.version=$(git rev-parse --short HEAD)"
var version = "dev"

// salud implementa adapterhttp.ProveedorSalud.
type salud struct {
	proveedorLLM string
	llmProvider  usecase.LLMProvider
	bd           string
	reloj        usecase.Reloj
}

func (s salud) Estado() adapterhttp.EstadoSalud {
	fecha := time.Now().UTC()
	if s.reloj != nil {
		fecha = s.reloj.FechaSimulada()
	}
	return adapterhttp.EstadoSalud{
		Estado:             "OK",
		Version:            version,
		LLMProveedorActivo: s.proveedorLLM,
		CircuitBreaker:     llm.CircuitBreakerHealth(s.llmProvider),
		BD:                 s.bd,
		FechaSimulada:      fecha.Format("2006-01-02"),
	}
}

type demoGateway struct{}

func (demoGateway) Enviar(context.Context, *domain.Mensaje) error { return nil }

func main() {
	ctx := context.Background()

	cfg, err := config.Cargar()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	provider, providerErr := llm.NewFromConfig(cfg)
	if providerErr != nil {
		log.Printf("llm: %v", providerErr)
	}

	pool, err := postgres.Conectar(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("base de datos: %v", err)
	}
	defer pool.Close()

	if err := postgres.Migrar(ctx, pool); err != nil {
		log.Fatalf("migraciones: %v", err)
	}
	log.Println("migraciones aplicadas")

	demoRepo := postgres.NuevoDemoRepository(pool)
	if err := (&usecase.CargarSeed{Repository: demoRepo, Habilitado: cfg.DemoSeed}).Ejecutar(ctx); err != nil {
		log.Fatalf("seed demo: %v", err)
	}
	reloj, err := infrareloj.NuevoPostgres(ctx, demoRepo)
	if err != nil {
		log.Fatalf("reloj simulado: %v", err)
	}

	mux := http.NewServeMux()
	providerName := "unconfigured"
	if provider != nil {
		providerName = provider.Nombre()
	}
	mux.HandleFunc("GET /salud", adapterhttp.HandlerSalud(salud{
		proveedorLLM: providerName,
		llmProvider:  provider,
		bd:           "OK",
		reloj:        reloj,
	}))

	// === BLOQUE A: registrar aquí repos, motor, LLM, casos de uso y rutas /api de conversación ===
	leadRepo := postgres.NuevoLeadRepository(pool)
	fichaRepo := postgres.NuevoFichaRepository(pool)
	catalogo, err := postgres.NuevoCatalogo(os.DirFS("data"))
	if err != nil {
		log.Fatalf("catálogo: %v", err)
	}
	busEventos := bus.Nuevo(slog.Default())
	ids := ids.NuevoGeneradorID()
	perfilador := &usecase.PerfilarLead{Leads: leadRepo, Catalogo: catalogo, IDs: ids, Bus: busEventos, Reloj: reloj}
	saludo := &usecase.SaludarLead{Leads: leadRepo, LLM: provider, IDs: ids, Reloj: reloj}
	planRepo := postgres.NuevoPlanRepository(pool)
	hitos := &usecase.EjecutarHitos{Leads: leadRepo, Planes: planRepo, Gateway: demoGateway{}, Reloj: reloj, IDs: ids, Bus: busEventos}
	agentes.Nueva(busEventos, agentes.Dependencias{LeadNuevo: saludo.Ejecutar, Nutricionista: hitos}).Registrar()
	procesarMensaje := &usecase.ProcesarMensaje{Leads: leadRepo, LLM: provider, IDs: ids, Bus: busEventos, Reloj: reloj, Saludo: saludo}
	turnos := adapterhttp.NuevoEjecutorTurnos(procesarMensaje, ids, reloj)
	defer turnos.Cerrar()
	controlador, err := adapterhttp.NuevoControlador(adapterhttp.Dependencias{Perfilar: perfilador, Leads: leadRepo, Fichas: fichaRepo, Planes: planRepo, Catalogo: catalogo, Turnos: turnos, Demo: demoRepo, Reloj: reloj, AvanzarDemo: &usecase.AvanzarDemo{Demo: demoRepo, Reloj: reloj, Bus: busEventos}, ReiniciarDemo: &usecase.ReiniciarDemo{Repository: demoRepo, Reloj: reloj, Habilitado: cfg.DemoSeed}})
	if err != nil {
		log.Fatalf("controlador HTTP: %v", err)
	}
	controlador.Registrar(mux)

	// (issues #13, #15, #16, #25)

	// === BLOQUE B: estáticos empaquetados y rutas del dashboard ===
	if err := adapterhttp.RegistrarEstaticos(mux); err != nil {
		log.Fatalf("estáticos: %v", err)
	}
	// (issues #29, #35)

	srv := &http.Server{
		Addr:              ":" + cfg.Puerto,
		Handler:           adapterhttp.NuevoLimitadorTasa(mux, cfg.TrustedProxyCIDRs),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("servidor escuchando en :%s", cfg.Puerto)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("servidor: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("apagando…")
	ctxApagado, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctxApagado)
}
