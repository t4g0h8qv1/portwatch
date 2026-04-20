package portwatch

import (
	"sync"
	"time"
)

// SuppressConfig holds configuration for the scan suppression manager.
type SuppressConfig struct {
	// DefaultTTL is how long a suppression lasts if no duration is specified.
	DefaultTTL time.Duration
}

// DefaultSuppressConfig returns a SuppressConfig with sensible defaults.
func DefaultSuppressConfig() SuppressConfig {
	return SuppressConfig{
		DefaultTTL: 30 * time.Minute,
	}
}

// suppressEntry holds the expiry time for a suppressed target.
type suppressEntry struct {
	expiry time.Time
}

// ScanSuppressManager tracks which targets have been suppressed from scanning.
type ScanSuppressManager struct {
	mu      sync.Mutex
	entries map[string]suppressEntry
	cfg     SuppressConfig
	now     func() time.Time
}

// NewScanSuppressManager creates a new ScanSuppressManager.
func NewScanSuppressManager(cfg SuppressConfig) (*ScanSuppressManager, error) {
	if cfg.DefaultTTL <= 0 {
		return nil, errInvalidSuppressTTL
	}
	return &ScanSuppressManager{
		entries: make(map[string]suppressEntry),
		cfg:     cfg,
		now:     time.Now,
	}, nil
}

// Suppress marks target as suppressed for the given duration.
// If d is zero, DefaultTTL is used.
func (s *ScanSuppressManager) Suppress(target string, d time.Duration) {
	if d <= 0 {
		d = s.cfg.DefaultTTL
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[target] = suppressEntry{expiry: s.now().Add(d)}
}

// IsSuppressed reports whether target is currently suppressed.
func (s *ScanSuppressManager) IsSuppressed(target string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[target]
	if !ok {
		return false
	}
	if s.now().After(e.expiry) {
		delete(s.entries, target)
		return false
	}
	return true
}

// Lift removes a suppression for target before it expires.
func (s *ScanSuppressManager) Lift(target string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, target)
}

// Targets returns all currently suppressed targets.
func (s *ScanSuppressManager) Targets() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.now()
	out := make([]string, 0, len(s.entries))
	for k, e := range s.entries {
		if now.Before(e.expiry) {
			out = append(out, k)
		} else {
			delete(s.entries, k)
		}
	}
	return out
}
