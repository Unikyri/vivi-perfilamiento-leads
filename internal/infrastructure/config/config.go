package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config contiene toda la configuración del proceso.
// Contrato v1.1 §8 — tabla de variables de entorno.
type Config struct {
	Puerto       string
	DatabaseURL  string
	LLMProvider  string
	GeminiAPIKey string
	QwenAPIKey   string
	QwenBaseURL  string
	LLMFallback  string
	TasaEA       float64
	DemoSeed     bool
	LogNivel     string
}

// Cargar lee la configuración del entorno y valida lo obligatorio.
func Cargar() (Config, error) {
	c := Config{
		Puerto:       valor("PORT", "8080"),
		DatabaseURL:  valor("DATABASE_URL", ""),
		LLMProvider:  valor("LLM_PROVIDER", "gemini"),
		GeminiAPIKey: valor("GEMINI_API_KEY", ""),
		QwenAPIKey:   valor("QWEN_API_KEY", ""),
		QwenBaseURL:  valor("QWEN_BASE_URL", ""),
		LLMFallback:  valor("LLM_FALLBACK", "qwen"),
		LogNivel:     valor("LOG_NIVEL", "info"),
	}

	tasa, err := strconv.ParseFloat(valor("TASA_EA", "0.107"), 64)
	if err != nil {
		return Config{}, fmt.Errorf("TASA_EA no es un número válido: %w", err)
	}
	c.TasaEA = tasa
	c.DemoSeed = valor("DEMO_SEED", "false") == "true"

	if c.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL es obligatoria")
	}
	return c, nil
}

func valor(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}
