package llm

import (
	"sync"
	"time"
)

const (
	// A short conversational burst avoids queueing a user's immediate turns while
	// the refill rate remains conservative for provider APIs without new config.
	defaultProviderRateLimitBurst           = 3
	defaultProviderRateLimitRefillPerSecond = 0.5
)

// RateLimiter gates one logical invocation before it reaches a provider.
type RateLimiter interface {
	Allow() bool
}

type RateLimiterFunc func() bool

func (f RateLimiterFunc) Allow() bool { return f == nil || f() }

// TokenBucket is a bounded, clock-injected token bucket. One token permits one
// logical provider invocation; retries inside an adapter remain one invocation.
type TokenBucket struct {
	mu              sync.Mutex
	now             Clock
	capacity        float64
	tokens          float64
	refillPerSecond float64
	last            time.Time
}

func NewTokenBucket(clock Clock, capacity int, refillPerSecond float64) *TokenBucket {
	if clock == nil {
		clock = wallClock{}
	}
	if capacity < 1 {
		capacity = 1
	}
	if refillPerSecond < 0 {
		refillPerSecond = 0
	}
	now := clock.Now()
	return &TokenBucket{
		now:             clock,
		capacity:        float64(capacity),
		tokens:          float64(capacity),
		refillPerSecond: refillPerSecond,
		last:            now,
	}
}

// NewRateLimiter is an expressive alias for callers wiring provider limits.
func NewRateLimiter(clock Clock, capacity int, refillPerSecond float64) *TokenBucket {
	return NewTokenBucket(clock, capacity, refillPerSecond)
}

func (b *TokenBucket) Allow() bool {
	if b == nil {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.now.Now()
	if now.After(b.last) && b.refillPerSecond > 0 {
		b.tokens += now.Sub(b.last).Seconds() * b.refillPerSecond
		if b.tokens > b.capacity {
			b.tokens = b.capacity
		}
	}
	if now.After(b.last) {
		b.last = now
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
