package llm

import (
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }
type wallClock struct{}

func (wallClock) Now() time.Time { return time.Now() }

type BreakerState string

const (
	BreakerClosed   BreakerState = "closed"
	BreakerOpen     BreakerState = "open"
	BreakerHalfOpen BreakerState = "half_open"
)

type BreakerSnapshot struct {
	State    BreakerState
	Failures int
}
type Breaker struct {
	mu       sync.Mutex
	now      Clock
	failures int
	opened   time.Time
	halfOpen bool
}

func NewBreaker(now Clock) *Breaker {
	if now == nil {
		now = wallClock{}
	}
	return &Breaker{now: now}
}
func (b *Breaker) Allow() error {
	if b == nil {
		return providerError(KindConfig, nil)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.opened.IsZero() {
		return nil
	}
	if b.now.Now().Sub(b.opened) < 60*time.Second || b.halfOpen {
		return providerError(KindCircuitOpen, nil)
	}
	b.halfOpen = true
	return nil
}
func (b *Breaker) Success() {
	if b == nil {
		return
	}
	b.mu.Lock()
	b.failures, b.opened, b.halfOpen = 0, time.Time{}, false
	b.mu.Unlock()
}
func (b *Breaker) Failure(err error) {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if !Eligible(err) {
		b.failures = 0
		return
	}
	b.failures++
	if b.halfOpen || b.failures >= 3 {
		b.opened, b.halfOpen = b.now.Now(), false
	}
}
func (b *Breaker) Snapshot() BreakerSnapshot {
	if b == nil {
		return BreakerSnapshot{State: BreakerClosed}
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	state := BreakerClosed
	if !b.opened.IsZero() {
		state = BreakerOpen
		if b.now.Now().Sub(b.opened) >= 60*time.Second {
			state = BreakerHalfOpen
		}
	}
	return BreakerSnapshot{State: state, Failures: b.failures}
}
