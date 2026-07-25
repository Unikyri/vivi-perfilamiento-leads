package llm

import (
	"context"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/infrastructure/config"
	"github.com/Unikyri/vivi-perfilamiento-leads/internal/usecase"
	"net/http"
	"strings"
	"testing"
	"time"
)

type fakeLLM struct {
	name              string
	err               error
	calls, audioCalls int
}

func (f *fakeLLM) Nombre() string { return f.name }
func (f *fakeLLM) GenerarTurno(context.Context, usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	f.calls++
	return usecase.SalidaTurno{}, f.err
}
func (f *fakeLLM) ProcesarAudio(context.Context, usecase.Audio, usecase.EntradaTurno) (usecase.SalidaTurno, error) {
	f.audioCalls++
	return usecase.SalidaTurno{}, f.err
}

type fakeClock struct{ now time.Time }

func (c *fakeClock) Now() time.Time      { return c.now }
func (c *fakeClock) add(d time.Duration) { c.now = c.now.Add(d) }
func fallbackFromDecorated(provider usecase.LLMProvider) *FallbackProvider {
	if metrics, ok := provider.(*Metricas); ok {
		provider = metrics.next
	}
	if guardrails, ok := provider.(*Guardarrailes); ok {
		provider = guardrails.next
	}
	composition, _ := provider.(*FallbackProvider)
	return composition
}

func TestFactorySelectionAndSafeIdentity(t *testing.T) {
	cfg := config.Config{LLMProvider: "qwen", QwenAPIKey: "qwen-secret"}
	p, err := NewFromConfig(cfg)
	if err != nil || p.Nombre() != "qwen" {
		t.Fatalf("provider=%v err=%v", p, err)
	}
	if got := HealthIdentity(cfg); got != "qwen" {
		t.Fatalf("identity=%q", got)
	}
	for _, tc := range []config.Config{{LLMProvider: "other"}, {LLMProvider: "gemini", GeminiAPIKey: ""}} {
		_, err = NewFromConfig(tc)
		if ErrorKindOf(err) != KindConfig || strings.Contains(err.Error(), "secret") {
			t.Fatalf("safe config error=%v", err)
		}
	}
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("LLM_PROVIDER", "gemini")
	t.Setenv("GEMINI_API_KEY", "gemini-key")
	t.Setenv("QWEN_API_KEY", "")
	for _, fallback := range []string{"", "qwen"} {
		t.Setenv("LLM_FALLBACK", fallback)
		cfg, loadErr := config.Cargar()
		if loadErr != nil {
			t.Fatal(loadErr)
		}
		p, err = NewFromConfig(cfg)
		if err != nil || p.Nombre() != "gemini" {
			t.Fatalf("fallback=%q provider=%v err=%v", fallback, p, err)
		}
	}
	if got := HealthIdentity(config.Config{LLMProvider: "gemini"}); got != "unconfigured" {
		t.Fatalf("identity=%q", got)
	}
}
func TestFallbackEligibilityAndExhaustion(t *testing.T) {
	primary := &fakeLLM{name: "gemini", err: providerError(KindRateLimit, nil)}
	backup := &fakeLLM{name: "qwen"}
	p := NewFallbackProvider(primary, backup)
	if _, err := p.GenerarTurno(context.Background(), usecase.EntradaTurno{}); err != nil || primary.calls != 1 || backup.calls != 1 {
		t.Fatalf("fallback err=%v calls=%d/%d", err, primary.calls, backup.calls)
	}
	for _, kind := range []ErrorKind{KindConfig, KindCanceled, KindHTTP, KindInvalidOutput, KindCapability} {
		primary = &fakeLLM{name: "gemini", err: providerError(kind, nil)}
		backup = &fakeLLM{name: "qwen"}
		_, err := NewFallbackProvider(primary, backup).GenerarTurno(context.Background(), usecase.EntradaTurno{})
		if ErrorKindOf(err) != kind || backup.calls != 0 {
			t.Fatalf("kind=%s err=%v backup=%d", kind, err, backup.calls)
		}
	}
	primary = &fakeLLM{name: "gemini", err: providerError(KindUnavailable, nil)}
	backup = &fakeLLM{name: "qwen", err: providerError(KindTimeout, nil)}
	_, err := NewFallbackProvider(primary, backup).GenerarTurno(context.Background(), usecase.EntradaTurno{})
	if ErrorKindOf(err) != KindExhausted || strings.Contains(err.Error(), "qwen") {
		t.Fatalf("exhaustion=%v", err)
	}
}
func TestQwenAudioNeverReroutes(t *testing.T) {
	qwen := NewQwenProvider("key", "http://unused", WithQwenHTTPDoer(fakeDoer(func(*http.Request) (*http.Response, error) { t.Fatal("unexpected HTTP"); return nil, nil })))
	backup := &fakeLLM{name: "gemini"}
	_, err := NewFallbackProvider(qwen, backup).ProcesarAudio(context.Background(), usecase.Audio{}, usecase.EntradaTurno{})
	if ErrorKindOf(err) != KindCapability || backup.audioCalls != 0 {
		t.Fatalf("err=%v backup=%d", err, backup.audioCalls)
	}
}
func TestBreakerOpenResetAndTimeout(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	b := NewBreaker(clock)
	for i := 0; i < 3; i++ {
		b.Failure(providerError(KindTimeout, nil))
	}
	if ErrorKindOf(b.Allow()) != KindCircuitOpen || b.Snapshot().State != BreakerOpen {
		t.Fatal("breaker did not open")
	}
	clock.add(59 * time.Second)
	if ErrorKindOf(b.Allow()) != KindCircuitOpen {
		t.Fatal("breaker opened too briefly")
	}
	clock.add(time.Second)
	if err := b.Allow(); err != nil || b.Snapshot().State != BreakerHalfOpen {
		t.Fatalf("probe err=%v", err)
	}
	b.Success()
	if got := b.Snapshot(); got.State != BreakerClosed || got.Failures != 0 {
		t.Fatalf("reset=%+v", got)
	}
}
func TestGuardrailStopsCallsWithoutLeakingInput(t *testing.T) {
	primary := &fakeLLM{name: "gemini"}
	backup := &fakeLLM{name: "qwen"}
	secret := "cedula-secret-message"
	p := NewFallbackProvider(primary, backup, WithGuardrail(GuardrailFunc(func(context.Context, usecase.EntradaTurno) error { return providerError(KindGuardrail, nil) })))
	_, err := p.GenerarTurno(context.Background(), usecase.EntradaTurno{MensajeUsuario: secret, LeadID: secret})
	if ErrorKindOf(err) != KindGuardrail || primary.calls != 0 || backup.calls != 0 || strings.Contains(err.Error(), secret) {
		t.Fatalf("guardrail err=%v calls=%d/%d", err, primary.calls, backup.calls)
	}
}

func TestMalformedAfterRepairFallsBack(t *testing.T) {
	primaryCalls := 0
	primary := NewGeminiProvider("key", WithGeminiHTTPDoer(fakeDoer(func(*http.Request) (*http.Response, error) {
		primaryCalls++
		return response(http.StatusOK, wireResponse("{")), nil
	})), WithGeminiBaseURL("http://primary"))
	backup := &fakeLLM{name: "qwen"}
	p := NewFallbackProvider(primary, backup)
	if _, err := p.GenerarTurno(context.Background(), adapterInput()); err != nil {
		t.Fatalf("fallback err=%v", err)
	}
	if primaryCalls != 2 || backup.calls != 1 {
		t.Fatalf("calls primary=%d fallback=%d", primaryCalls, backup.calls)
	}
}

func TestTokenBucketDenialUsesFallback(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	primary := &fakeLLM{name: "gemini"}
	backup := &fakeLLM{name: "qwen"}
	p := NewFallbackProvider(primary, backup,
		WithPrimaryRateLimiter(NewTokenBucket(clock, 1, 1)),
		WithFallbackRateLimiter(RateLimiterFunc(func() bool { return true })),
	)
	if _, err := p.GenerarTurno(context.Background(), usecase.EntradaTurno{}); err != nil || primary.calls != 1 || backup.calls != 0 {
		t.Fatalf("first err=%v calls=%d/%d", err, primary.calls, backup.calls)
	}
	if _, err := p.GenerarTurno(context.Background(), usecase.EntradaTurno{}); err != nil || primary.calls != 1 || backup.calls != 1 {
		t.Fatalf("denied err=%v calls=%d/%d", err, primary.calls, backup.calls)
	}
	clock.add(time.Second)
	if _, err := p.GenerarTurno(context.Background(), usecase.EntradaTurno{}); err != nil || primary.calls != 2 {
		t.Fatalf("refill err=%v calls=%d", err, primary.calls)
	}
}

func TestHealthReflectsOpenPrimaryBreaker(t *testing.T) {
	primary := &fakeLLM{name: "gemini", err: providerError(KindTimeout, nil)}
	p := NewFallbackProvider(primary, nil, WithPrimaryRateLimiter(RateLimiterFunc(func() bool { return true })))
	for i := 0; i < 3; i++ {
		_, _ = p.GenerarTurno(context.Background(), usecase.EntradaTurno{})
	}
	if got := CircuitBreakerHealth(p); got != "ABIERTO" {
		t.Fatalf("health=%q", got)
	}
}

func TestFactoryNoFallbackRetainsBreakerAndLiveHealth(t *testing.T) {
	provider, err := NewFromConfig(config.Config{LLMProvider: "gemini", GeminiAPIKey: "key"})
	if err != nil {
		t.Fatal(err)
	}
	metrics, ok := provider.(*Metricas)
	if !ok {
		t.Fatalf("provider=%T; expected metrics decorator", provider)
	}
	if _, ok := metrics.next.(*Guardarrailes); !ok {
		t.Fatalf("inner=%T; expected guardrails decorator", metrics.next)
	}
	composition := fallbackFromDecorated(provider)
	if composition == nil {
		t.Fatalf("decorated provider did not preserve fallback composite")
	}
	if composition.fallback != nil {
		t.Fatalf("fallback=%v; expected no fallback", composition.fallback)
	}
	composition.primary = &fakeLLM{name: "gemini", err: providerError(KindTimeout, nil)}
	composition.primaryLimiter = RateLimiterFunc(func() bool { return true })
	for i := 0; i < 3; i++ {
		if _, callErr := composition.GenerarTurno(context.Background(), usecase.EntradaTurno{}); ErrorKindOf(callErr) != KindTimeout {
			t.Fatalf("failure %d kind=%q err=%v", i+1, ErrorKindOf(callErr), callErr)
		}
	}
	if got := CircuitBreakerHealth(provider); got != "ABIERTO" {
		t.Fatalf("live health=%q", got)
	}
}

func TestFactoryNoFallbackRetainsRateLimiter(t *testing.T) {
	provider, err := NewFromConfig(config.Config{LLMProvider: "gemini", GeminiAPIKey: "key", LLMFallback: ""})
	if err != nil {
		t.Fatal(err)
	}
	composition := fallbackFromDecorated(provider)
	if composition == nil {
		t.Fatalf("decorated provider did not preserve fallback composite")
	}
	primary := &fakeLLM{name: "gemini"}
	composition.primary = primary
	composition.primaryLimiter = NewTokenBucket(&fakeClock{now: time.Unix(0, 0)}, 1, 0)
	if _, err := composition.GenerarTurno(context.Background(), usecase.EntradaTurno{}); err != nil {
		t.Fatalf("first call err=%v", err)
	}
	if _, err := composition.GenerarTurno(context.Background(), usecase.EntradaTurno{}); ErrorKindOf(err) != KindRateLimit || primary.calls != 1 {
		t.Fatalf("limiter err=%v calls=%d", err, primary.calls)
	}
}

func TestBreakerRequiresThreeConsecutiveEligibleFailures(t *testing.T) {
	cases := []struct {
		name  string
		reset func(*Breaker)
	}{
		{name: "ineligible invocation", reset: func(b *Breaker) { b.Failure(providerError(KindInvalidOutput, nil)) }},
		{name: "successful invocation", reset: func(b *Breaker) { b.Success() }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewBreaker(&fakeClock{now: time.Unix(0, 0)})
			b.Failure(providerError(KindTimeout, nil))
			tc.reset(b)
			b.Failure(providerError(KindTimeout, nil))
			b.Failure(providerError(KindTimeout, nil))
			if err := b.Allow(); err != nil {
				t.Fatalf("breaker opened before three consecutive failures: %v", err)
			}
			b.Failure(providerError(KindTimeout, nil))
			if ErrorKindOf(b.Allow()) != KindCircuitOpen {
				t.Fatal("breaker did not open on the third consecutive eligible failure")
			}
		})
	}
}

func TestDefaultProviderRateLimiterUsesNamedBoundedChatDefaults(t *testing.T) {
	clock := &fakeClock{now: time.Unix(0, 0)}
	composition := NewFallbackProvider(&fakeLLM{name: "gemini"}, &fakeLLM{name: "qwen"}, WithFallbackClock(clock))
	limiter, ok := composition.primaryLimiter.(*TokenBucket)
	if !ok {
		t.Fatalf("primary limiter=%T; want bounded token bucket", composition.primaryLimiter)
	}
	if limiter.capacity != defaultProviderRateLimitBurst || limiter.refillPerSecond != defaultProviderRateLimitRefillPerSecond {
		t.Fatalf("defaults capacity=%v refill=%v", limiter.capacity, limiter.refillPerSecond)
	}
	for i := 0; i < defaultProviderRateLimitBurst; i++ {
		if !limiter.Allow() {
			t.Fatalf("default burst token %d denied", i+1)
		}
	}
	if limiter.Allow() {
		t.Fatal("default limiter exceeded its bounded burst")
	}
	clock.add(2 * time.Second)
	if !limiter.Allow() {
		t.Fatal("default limiter did not refill at its documented rate")
	}
}

func TestCompositeBreakerBypassesPrimaryAfterThreeEligibleFailures(t *testing.T) {
	primary := &fakeLLM{name: "gemini", err: providerError(KindTimeout, nil)}
	fallback := &fakeLLM{name: "qwen"}
	provider := NewFallbackProvider(primary, fallback,
		WithPrimaryRateLimiter(RateLimiterFunc(func() bool { return true })),
		WithFallbackRateLimiter(RateLimiterFunc(func() bool { return true })),
	)
	for i := 0; i < 4; i++ {
		if _, err := provider.GenerarTurno(context.Background(), usecase.EntradaTurno{}); err != nil {
			t.Fatalf("invocation %d err=%v", i+1, err)
		}
	}
	if primary.calls != 3 {
		t.Fatalf("primary calls=%d; want 3 before fourth invocation bypass", primary.calls)
	}
	if fallback.calls != 4 {
		t.Fatalf("fallback calls=%d; want all four invocations", fallback.calls)
	}
	if got := provider.CircuitBreakerState(); got != BreakerOpen {
		t.Fatalf("breaker state=%q; want %q", got, BreakerOpen)
	}
}
