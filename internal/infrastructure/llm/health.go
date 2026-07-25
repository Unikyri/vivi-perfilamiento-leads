package llm

import "github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"

// CircuitBreakerHealth returns the stable health vocabulary used by /salud.
// Non-composite providers are closed because they do not own a breaker.
func CircuitBreakerHealth(provider usecase.LLMProvider) string {
	if provider == nil || provider.Nombre() == "unconfigured" {
		return "CERRADO"
	}
	state := BreakerClosed
	if composite, ok := provider.(*FallbackProvider); ok {
		state = composite.CircuitBreakerState()
	}
	switch state {
	case BreakerOpen:
		return "ABIERTO"
	case BreakerHalfOpen:
		return "MEDIO_ABIERTO"
	default:
		return "CERRADO"
	}
}
