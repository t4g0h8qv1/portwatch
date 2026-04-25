package portwatch

import (
	"errors"
	"sync"
	"time"
)

// DefaultProbeConfig returns a ProbeConfig with sensible defaults.
func DefaultProbeConfig() ProbeConfig {
	return ProbeConfig{
		MaxConsecutiveFailures: 3,
		ProbeInterval:          30 * time.Second,
		RecoveryThreshold:      2,
	}
}

// ProbeConfig controls how the probe manager tracks target liveness.
type ProbeConfig struct {
	MaxConsecutiveFailures int
	ProbeInterval          time.Duration
	RecoveryThreshold      int
}

// probeState holds per-target probe tracking state.
type probeState struct {
	consecFailures  int
	consecSuccesses int
	lastProbe       time.Time
	dead            bool
}

// ProbeManager tracks consecutive scan failures per target and marks
// targets as dead once the failure threshold is exceeded. A target
// recovers after RecoveryThreshold consecutive successes.
type ProbeManager struct {
	mu     sync.Mutex
	cfg    ProbeConfig
	states map[string]*probeState
}

// NewProbeManager creates a ProbeManager with the given config.
func NewProbeManager(cfg ProbeConfig) (*ProbeManager, error) {
	if cfg.MaxConsecutiveFailures < 1 {
		return nil, errors.New("portwatch: MaxConsecutiveFailures must be >= 1")
	}
	if cfg.ProbeInterval <= 0 {
		return nil, errors.New("portwatch: ProbeInterval must be positive")
	}
	if cfg.RecoveryThreshold < 1 {
		return nil, errors.New("portwatch: RecoveryThreshold must be >= 1")
	}
	return &ProbeManager{
		cfg:    cfg,
		states: make(map[string]*probeState),
	}, nil
}

func (m *ProbeManager) state(target string) *probeState {
	s, ok := m.states[target]
	if !ok {
		s = &probeState{}
		m.states[target] = s
	}
	return s
}

// RecordSuccess records a successful scan for target.
func (m *ProbeManager) RecordSuccess(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := m.state(target)
	s.consecFailures = 0
	s.lastProbe = time.Now()
	if s.dead {
		s.consecSuccesses++
		if s.consecSuccesses >= m.cfg.RecoveryThreshold {
			s.dead = false
			s.consecSuccesses = 0
		}
	}
}

// RecordFailure records a failed scan for target.
func (m *ProbeManager) RecordFailure(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s := m.state(target)
	s.consecSuccesses = 0
	s.lastProbe = time.Now()
	s.consecFailures++
	if s.consecFailures >= m.cfg.MaxConsecutiveFailures {
		s.dead = true
	}
}

// IsDead reports whether target has exceeded the failure threshold.
func (m *ProbeManager) IsDead(target string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state(target).dead
}

// Reset clears probe state for target.
func (m *ProbeManager) Reset(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.states, target)
}
