package portwatch

import (
	"fmt"
	"sync"
	"time"
)

// StalenessConfig holds configuration for the stale-scan detector.
type StalenessConfig struct {
	// MaxAge is the duration after which a scan result is considered stale.
	MaxAge time.Duration
}

// DefaultStalenessConfig returns a StalenessConfig with sensible defaults.
func DefaultStalenessConfig() StalenessConfig {
	return StalenessConfig{
		MaxAge: 10 * time.Minute,
	}
}

// StalenessManager tracks the last successful scan time per target and
// reports whether the most recent result has exceeded the configured MaxAge.
type StalenessManager struct {
	cfg  StalenessConfig
	mu   sync.RWMutex
	last map[string]time.Time
	now  func() time.Time
}

// NewStalenessManager creates a StalenessManager with the given config.
// Returns an error if MaxAge is not positive.
func NewStalenessManager(cfg StalenessConfig) (*StalenessManager, error) {
	if cfg.MaxAge <= 0 {
		return nil, fmt.Errorf("portwatch: stale max age must be positive, got %v", cfg.MaxAge)
	}
	return &StalenessManager{
		cfg:  cfg,
		last: make(map[string]time.Time),
		now:  time.Now,
	}, nil
}

// Observe records the time of the most recent successful scan for target.
func (s *StalenessManager) Observe(target string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last[target] = s.now()
}

// IsStale reports whether the last recorded scan for target occurred more
// than MaxAge ago, or if no scan has ever been recorded.
func (s *StalenessManager) IsStale(target string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.last[target]
	if !ok {
		return true
	}
	return s.now().Sub(t) > s.cfg.MaxAge
}

// LastScan returns the time of the most recent scan for target and whether
// one has been recorded.
func (s *StalenessManager) LastScan(target string) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.last[target]
	return t, ok
}

// Targets returns all targets that have been observed.
func (s *StalenessManager) Targets() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, 0, len(s.last))
	for k := range s.last {
		out = append(out, k)
	}
	return out
}
