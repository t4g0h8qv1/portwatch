package portwatch

import (
	"fmt"
	"sync"
	"time"
)

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "closed"
	case CircuitOpen:
		return "open"
	case CircuitHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig holds configuration for the circuit breaker.
type CircuitBreakerConfig struct {
	MaxFailures int
	OpenDuration time.Duration
}

// DefaultCircuitBreakerConfig returns sensible defaults.
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		MaxFailures:  3,
		OpenDuration: 30 * time.Second,
	}
}

type circuitEntry struct {
	failures  int
	state     CircuitState
	openUntil time.Time
}

// CircuitBreaker tracks per-target failure counts and opens/closes circuits.
type CircuitBreaker struct {
	mu     sync.Mutex
	cfg    CircuitBreakerConfig
	targets map[string]*circuitEntry
}

// NewCircuitBreaker creates a CircuitBreaker with the given config.
func NewCircuitBreaker(cfg CircuitBreakerConfig) (*CircuitBreaker, error) {
	if cfg.MaxFailures < 1 {
		return nil, fmt.Errorf("portwatch: MaxFailures must be >= 1")
	}
	if cfg.OpenDuration <= 0 {
		return nil, fmt.Errorf("portwatch: OpenDuration must be positive")
	}
	return &CircuitBreaker{
		cfg:     cfg,
		targets: make(map[string]*circuitEntry),
	}, nil
}

func (cb *CircuitBreaker) entry(target string) *circuitEntry {
	e, ok := cb.targets[target]
	if !ok {
		e = &circuitEntry{state: CircuitClosed}
		cb.targets[target] = e
	}
	return e
}

// Allow returns true if the target's circuit permits a scan attempt.
func (cb *CircuitBreaker) Allow(target string) bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	e := cb.entry(target)
	if e.state == CircuitOpen {
		if time.Now().After(e.openUntil) {
			e.state = CircuitHalfOpen
			return true
		}
		return false
	}
	return true
}

// RecordSuccess resets the failure count and closes the circuit.
func (cb *CircuitBreaker) RecordSuccess(target string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	e := cb.entry(target)
	e.failures = 0
	e.state = CircuitClosed
}

// RecordFailure increments the failure count and may open the circuit.
func (cb *CircuitBreaker) RecordFailure(target string) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	e := cb.entry(target)
	e.failures++
	if e.failures >= cb.cfg.MaxFailures {
		e.state = CircuitOpen
		e.openUntil = time.Now().Add(cb.cfg.OpenDuration)
	}
}

// State returns the current CircuitState for a target.
func (cb *CircuitBreaker) State(target string) CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.entry(target).state
}
