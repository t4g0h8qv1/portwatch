package portwatch

import (
	"errors"
	"sync"
	"time"
)

// DefaultPressureConfig returns a PressureConfig with sensible defaults.
func DefaultPressureConfig() PressureConfig {
	return PressureConfig{
		HighWatermark: 0.80,
		LowWatermark:  0.50,
		Window:        30 * time.Second,
	}
}

// PressureConfig controls thresholds for scan pressure management.
type PressureConfig struct {
	// HighWatermark is the fraction of capacity (0–1) at which pressure is
	// considered high and new scans should be shed.
	HighWatermark float64
	// LowWatermark is the fraction below which pressure is considered normal.
	LowWatermark float64
	// Window is the rolling window used to measure scan rate.
	Window time.Duration
}

// PressureLevel describes the current load on the scanner.
type PressureLevel int

const (
	PressureNormal PressureLevel = iota
	PressureHigh
)

// String returns a human-readable label for the level.
func (p PressureLevel) String() string {
	switch p {
	case PressureHigh:
		return "high"
	default:
		return "normal"
	}
}

// ScanPressureManager tracks scan throughput and reports pressure levels.
type ScanPressureManager struct {
	cfg    PressureConfig
	max    int
	mu     sync.Mutex
	times  []time.Time
	nowFn  func() time.Time
}

// NewScanPressureManager creates a manager that tracks pressure against a
// maximum expected scans-per-window capacity.
func NewScanPressureManager(cfg PressureConfig, maxPerWindow int) (*ScanPressureManager, error) {
	if maxPerWindow <= 0 {
		return nil, errors.New("portwatch: maxPerWindow must be positive")
	}
	if cfg.HighWatermark <= 0 || cfg.HighWatermark > 1 {
		return nil, errors.New("portwatch: HighWatermark must be in (0,1]")
	}
	if cfg.LowWatermark <= 0 || cfg.LowWatermark >= cfg.HighWatermark {
		return nil, errors.New("portwatch: LowWatermark must be in (0, HighWatermark)")
	}
	if cfg.Window <= 0 {
		return nil, errors.New("portwatch: Window must be positive")
	}
	return &ScanPressureManager{
		cfg:   cfg,
		max:   maxPerWindow,
		nowFn: time.Now,
	}, nil
}

// Record registers that a scan occurred now.
func (m *ScanPressureManager) Record() {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.nowFn()
	m.prune(now)
	m.times = append(m.times, now)
}

// Level returns the current PressureLevel based on recent scan volume.
func (m *ScanPressureManager) Level() PressureLevel {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.nowFn()
	m.prune(now)
	ratio := float64(len(m.times)) / float64(m.max)
	if ratio >= m.cfg.HighWatermark {
		return PressureHigh
	}
	return PressureNormal
}

// Load returns the current utilisation fraction (0–1).
func (m *ScanPressureManager) Load() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.nowFn()
	m.prune(now)
	return float64(len(m.times)) / float64(m.max)
}

// prune removes observations older than the configured window. Must be called
// with m.mu held.
func (m *ScanPressureManager) prune(now time.Time) {
	cutoff := now.Add(-m.cfg.Window)
	i := 0
	for i < len(m.times) && m.times[i].Before(cutoff) {
		i++
	}
	m.times = m.times[i:]
}
