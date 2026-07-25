package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	adapterhttp "github.com/Unikyri/vivi-perfilamiento-leads/internal/adapters/http"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/config"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/llm"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/postgres"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

// version la inyecta el build: -ldflags "-X main.version=$(git rev-parse --short HEAD)"
var version = "dev"

// salud implementa adapterhttp.ProveedorSalud.
type salud struct {
	proveedorLLM string
	llmProvider  usecase.LLMProvider
	bd           string
}

func (s salud) Estado() adapterhttp.EstadoSalud {
	return adapterhttp.EstadoSalud{
		Estado:             "OK",
		Version:            version,
		LLMProveedorActivo: s.proveedorLLM,
		CircuitBreaker:     llm.CircuitBreakerHealth(s.llmProvider),
		BD:                 s.bd,
		FechaSimulada:      time.Now().Format("2006-01-02"),
	}
}

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

	mux := http.NewServeMux()
	providerName := "unconfigured"
	if provider != nil {
		providerName = provider.Nombre()
	}
	mux.HandleFunc("GET /salud", adapterhttp.HandlerSalud(salud{
		proveedorLLM: providerName,
		llmProvider:  provider,
		bd:           "OK",
	}))

	// === BLOQUE A: registrar aquí repos, motor, LLM, casos de uso y rutas /api de conversación ===
	// (issues #13, #15, #16, #25)

	// === BLOQUE B: registrar aquí el servido de estáticos de web/ y las rutas del dashboard ===
	// (issues #29, #35)

	srv := &http.Server{
		Addr:              ":" + cfg.Puerto,
		Handler:           mux,
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
