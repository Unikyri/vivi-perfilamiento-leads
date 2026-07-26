package llm

import (
	"context"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"time"
)

const LogicalTimeout = 10 * time.Second

type Guardrail interface {
	Check(context.Context, usecase.EntradaTurno) error
}
type GuardrailFunc func(context.Context, usecase.EntradaTurno) error

func (f GuardrailFunc) Check(ctx context.Context, in usecase.EntradaTurno) error { return f(ctx, in) }

type Metrics interface{ Observe(string, ErrorKind) }
type NopMetrics struct{}

func (NopMetrics) Observe(string, ErrorKind) {}

type FallbackOption func(*FallbackProvider)

func WithGuardrail(g Guardrail) FallbackOption { return func(p *FallbackProvider) { p.guardrail = g } }
func WithMetrics(m Metrics) FallbackOption {
	return func(p *FallbackProvider) {
		if m != nil {
			p.metrics = m
		}
	}
}
func WithFallbackClock(c Clock) FallbackOption {
	return func(p *FallbackProvider) {
		if c != nil {
			p.clock = c
		}
	}
}
func WithPrimaryRateLimiter(l RateLimiter) FallbackOption {
	return func(p *FallbackProvider) { p.primaryLimiter = l }
}
func WithFallbackRateLimiter(l RateLimiter) FallbackOption {
	return func(p *FallbackProvider) { p.fallbackLimiter = l }
}
func WithPrimaryBreaker(b *Breaker) FallbackOption {
	return func(p *FallbackProvider) { p.primaryBreaker, p.primaryCustom = b, true }
}
func WithFallbackBreaker(b *Breaker) FallbackOption {
	return func(p *FallbackProvider) { p.fallbackBreaker, p.fallbackCustom = b, true }
}

type FallbackProvider struct {
	primary, fallback               usecase.LLMProvider
	primaryBreaker, fallbackBreaker *Breaker
	primaryCustom, fallbackCustom   bool
	primaryLimiter, fallbackLimiter RateLimiter
	clock                           Clock
	guardrail                       Guardrail
	metrics                         Metrics
}

func NewFallbackProvider(primary, fallback usecase.LLMProvider, options ...FallbackOption) *FallbackProvider {
	p := &FallbackProvider{primary: primary, fallback: fallback, clock: wallClock{}, metrics: NopMetrics{}}
	for _, option := range options {
		option(p)
	}
	if !p.primaryCustom {
		p.primaryBreaker = NewBreaker(p.clock)
	}
	if !p.fallbackCustom {
		p.fallbackBreaker = NewBreaker(p.clock)
	}
	if p.primaryLimiter == nil {
		p.primaryLimiter = NewTokenBucket(p.clock, defaultProviderRateLimitBurst, defaultProviderRateLimitRefillPerSecond)
	}
	if p.fallbackLimiter == nil {
		p.fallbackLimiter = NewTokenBucket(p.clock, defaultProviderRateLimitBurst, defaultProviderRateLimitRefillPerSecond)
	}
	return p
}
func (p *FallbackProvider) Nombre() string {
	if p == nil || p.primary == nil {
		return "unconfigured"
	}
	return p.primary.Nombre()
}

// CircuitBreakerState exposes the live primary breaker used by this composition.
func (p *FallbackProvider) CircuitBreakerState() BreakerState {
	if p == nil || p.primaryBreaker == nil {
		return BreakerClosed
	}
	return p.primaryBreaker.Snapshot().State
}

func (p *FallbackProvider) GenerarTurno(ctx context.Context, in usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	return p.run(ctx, in, func(c context.Context, provider usecase.LLMProvider) (usecase.SalidaTurno, error) {
		return provider.GenerarTurno(c, in)
	})
}
func (p *FallbackProvider) ProcesarAudio(ctx context.Context, audio usecase.Audio, in usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	return p.run(ctx, in, func(c context.Context, provider usecase.LLMProvider) (usecase.SalidaTurno, error) {
		return provider.ProcesarAudio(c, audio, in)
	})
}
func (p *FallbackProvider) run(ctx context.Context, in usecase.EntradaTurno, call func(context.Context, usecase.LLMProvider) (usecase.SalidaTurno, error)) (usecase.SalidaTurno, error) {
	if p == nil || p.primary == nil {
		return usecase.SalidaTurno{}, providerError(KindConfig, nil)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if p.guardrail != nil {
		if err := p.guardrail.Check(ctx, in); err != nil {
			if ErrorKindOf(err) != "" {
				return usecase.SalidaTurno{}, err
			}
			return usecase.SalidaTurno{}, providerError(KindGuardrail, nil)
		}
	}
	op, cancel := context.WithTimeout(ctx, LogicalTimeout)
	defer cancel()
	out, err := p.try(op, p.primary, p.primaryBreaker, p.primaryLimiter, call)
	if err == nil || !Eligible(err) || op.Err() != nil || p.fallback == nil {
		return out, err
	}
	backup, backupErr := p.try(op, p.fallback, p.fallbackBreaker, p.fallbackLimiter, call)
	if backupErr == nil {
		return backup, nil
	}
	if !Eligible(backupErr) {
		return backup, backupErr
	}
	return usecase.SalidaTurno{}, providerError(KindExhausted, nil)
}
func (p *FallbackProvider) try(ctx context.Context, provider usecase.LLMProvider, breaker *Breaker, limiter RateLimiter, call func(context.Context, usecase.LLMProvider) (usecase.SalidaTurno, error)) (usecase.SalidaTurno, error) {
	if provider == nil || breaker == nil {
		return usecase.SalidaTurno{}, providerError(KindConfig, nil)
	}
	if err := breaker.Allow(); err != nil {
		p.metrics.Observe(provider.Nombre(), ErrorKindOf(err))
		return usecase.SalidaTurno{}, err
	}
	if limiter != nil && !limiter.Allow() {
		err := providerError(KindRateLimit, nil)
		p.metrics.Observe(provider.Nombre(), ErrorKindOf(err))
		return usecase.SalidaTurno{}, err
	}
	out, err := call(ctx, provider)
	if err == nil {
		breaker.Success()
	} else {
		breaker.Failure(err)
	}
	p.metrics.Observe(provider.Nombre(), ErrorKindOf(err))
	return out, err
}
