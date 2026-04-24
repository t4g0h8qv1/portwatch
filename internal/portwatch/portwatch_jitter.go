package portwatch

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// JitterConfig holds configuration for scan jitter.
type JitterConfig struct {
	// MaxJitter is the upper bound on the random delay added before each scan.
	MaxJitter time.Duration
}

// DefaultJitterConfig returns a JitterConfig with sensible defaults.
func DefaultJitterConfig() JitterConfig {
	return JitterConfig{
		MaxJitter: 5 * time.Second,
	}
}

// JitterManager adds a random delay before scans to spread load across targets.
type JitterManager struct {
	cfg  JitterConfig
	mu   sync.Mutex
	rng  *rand.Rand
}

// NewJitterManager creates a JitterManager with the given config.
func NewJitterManager(cfg JitterConfig) (*JitterManager, error) {
	if cfg.MaxJitter <= 0 {
		return nil, errors.New("portwatch: jitter MaxJitter must be greater than zero")
	}
	return &JitterManager{
		cfg: cfg,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Delay returns a random duration in [0, MaxJitter).
func (j *JitterManager) Delay() time.Duration {
	j.mu.Lock()
	defer j.mu.Unlock()
	n := j.rng.Int63n(int64(j.cfg.MaxJitter))
	return time.Duration(n)
}

// Wait blocks for a random jitter duration.
func (j *JitterManager) Wait() {
	time.Sleep(j.Delay())
}

// MaxJitter returns the configured maximum jitter.
func (j *JitterManager) MaxJitter() time.Duration {
	return j.cfg.MaxJitter
}
