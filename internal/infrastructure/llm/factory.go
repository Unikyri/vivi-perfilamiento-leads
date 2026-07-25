package llm

import (
	"strings"

	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/config"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
)

func NewFromConfig(cfg config.Config) (usecase.LLMProvider, error) {
	primary, err := providerFor(cfg.LLMProvider, cfg)
	if err != nil {
		return nil, err
	}
	fallbackName := strings.TrimSpace(strings.ToLower(cfg.LLMFallback))
	if fallbackName == "" || fallbackName == strings.ToLower(primary.Nombre()) {
		return NewFallbackProvider(primary, nil), nil
	}
	fallback, _ := providerFor(fallbackName, cfg)
	return NewFallbackProvider(primary, fallback), nil
}
func providerFor(name string, cfg config.Config) (usecase.LLMProvider, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "gemini":
		if strings.TrimSpace(cfg.GeminiAPIKey) == "" {
			return nil, providerError(KindConfig, nil)
		}
		return NewGeminiProvider(cfg.GeminiAPIKey), nil
	case "qwen":
		if strings.TrimSpace(cfg.QwenAPIKey) == "" {
			return nil, providerError(KindConfig, nil)
		}
		return NewQwenProvider(cfg.QwenAPIKey, cfg.QwenBaseURL), nil
	default:
		return nil, providerError(KindConfig, nil)
	}
}

func HealthIdentity(cfg config.Config) string {
	p, err := NewFromConfig(cfg)
	if err != nil {
		return "unconfigured"
	}
	return p.Nombre()
}
