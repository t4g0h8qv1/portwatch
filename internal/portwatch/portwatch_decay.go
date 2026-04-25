package portwatch

import (
	"errors"
	"sync"
	"time"
)

// DecayConfig holds configuration for the scan decay manager.
type DecayConfig struct {
	// HalfLife is the duration after which a score decays to half its value.
	HalfLife time.Duration
	// InitialScore is the score assigned when a target is first observed.
	InitialScore float64
}

// DefaultDecayConfig returns a DecayConfig with sensible defaults.
func DefaultDecayConfig() DecayConfig {
	return DecayConfig{
		HalfLife:     30 * time.Minute,
		InitialScore: 100.0,
	}
}

type decayEntry struct {
	score     float64
	updatedAt time.Time
}

// DecayManager tracks a decaying relevance score per scan target.
// Scores decay exponentially based on elapsed time and the configured half-life.
type DecayManager struct {
	mu      sync.Mutex
	cfg     DecayConfig
	entries map[string]decayEntry
	now     func() time.Time
}

// NewDecayManager creates a DecayManager with the given config.
func NewDecayManager(cfg DecayConfig) (*DecayManager, error) {
	if cfg.HalfLife <= 0 {
		return nil, errors.New("portwatch: decay half-life must be positive")
	}
	if cfg.InitialScore <= 0 {
		return nil, errors.New("portwatch: decay initial score must be positive")
	}
	return &DecayManager{
		cfg:     cfg,
		entries: make(map[string]decayEntry),
		now:     time.Now,
	}, nil
}

// Observe records a fresh observation for target, resetting its score to InitialScore.
func (d *DecayManager) Observe(target string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries[target] = decayEntry{
		score:     d.cfg.InitialScore,
		updatedAt: d.now(),
	}
}

// Score returns the current decayed score for target.
// Returns 0 if the target has never been observed.
func (d *DecayManager) Score(target string) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	e, ok := d.entries[target]
	if !ok {
		return 0
	}
	elapsed := d.now().Sub(e.updatedAt)
	// exponential decay: score * 0.5^(elapsed/halfLife)
	exponent := float64(elapsed) / float64(d.cfg.HalfLife)
	decayed := e.score * pow2neg(exponent)
	return decayed
}

// Reset removes all decay state for the given target.
func (d *DecayManager) Reset(target string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.entries, target)
}

// Targets returns all tracked targets.
func (d *DecayManager) Targets() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]string, 0, len(d.entries))
	for k := range d.entries {
		out = append(out, k)
	}
	return out
}

// pow2neg computes 2^(-exp) using natural exponentiation.
func pow2neg(exp float64) float64 {
	// 2^(-x) = e^(-x * ln2)
	const ln2 = 0.6931471805599453
	return expApprox(-exp * ln2)
}

// expApprox is a small wrapper so we don't import math in the main logic.
func expApprox(x float64) float64 {
	// Use the standard library via a thin indirection captured at init.
	return mathExp(x)
}

var mathExp = func(x float64) float64 {
	// real implementation delegates to math.Exp
	import_math_exp := _mathExp
	return import_math_exp(x)
}

func init() {
	mathExp = _mathExp
}
