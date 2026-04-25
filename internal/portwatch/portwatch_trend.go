package portwatch

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// DefaultTrendConfig returns a TrendConfig with sensible defaults.
func DefaultTrendConfig() TrendConfig {
	return TrendConfig{
		WindowSize: 10,
		MinSamples: 3,
		RiseThreshold: 0.5,
		FallThreshold: 0.5,
	}
}

// TrendConfig controls how port-open trends are evaluated.
type TrendConfig struct {
	WindowSize    int     // number of recent scans to consider
	MinSamples    int     // minimum samples before a trend is reported
	RiseThreshold float64 // fraction of window that must be open to flag a rise
	FallThreshold float64 // fraction of window that must be closed to flag a fall
}

// TrendDirection indicates whether a port is trending open, closed, or stable.
type TrendDirection int

const (
	TrendStable TrendDirection = iota
	TrendRising
	TrendFalling
)

func (d TrendDirection) String() string {
	switch d {
	case TrendRising:
		return "rising"
	case TrendFalling:
		return "falling"
	default:
		return "stable"
	}
}

// PortTrend holds trend information for a single port on a target.
type PortTrend struct {
	Target    string
	Port      int
	Direction TrendDirection
	OpenRate  float64
	Samples   int
	UpdatedAt time.Time
}

// ScanTrendManager tracks open/closed observations per port per target
// and derives trend directions over a sliding window.
type ScanTrendManager struct {
	mu      sync.Mutex
	cfg     TrendConfig
	// observations[target][port] = []bool (true=open)
	obs     map[string]map[int][]bool
}

// NewScanTrendManager creates a ScanTrendManager with the given config.
func NewScanTrendManager(cfg TrendConfig) (*ScanTrendManager, error) {
	if cfg.WindowSize < 1 {
		return nil, fmt.Errorf("portwatch: TrendConfig.WindowSize must be >= 1")
	}
	if cfg.MinSamples < 1 {
		return nil, fmt.Errorf("portwatch: TrendConfig.MinSamples must be >= 1")
	}
	if cfg.RiseThreshold <= 0 || cfg.RiseThreshold > 1 {
		return nil, fmt.Errorf("portwatch: TrendConfig.RiseThreshold must be in (0,1]")
	}
	if cfg.FallThreshold <= 0 || cfg.FallThreshold > 1 {
		return nil, fmt.Errorf("portwatch: TrendConfig.FallThreshold must be in (0,1]")
	}
	return &ScanTrendManager{
		cfg: cfg,
		obs: make(map[string]map[int][]bool),
	}, nil
}

// Observe records whether a port was open during a scan for the given target.
func (m *ScanTrendManager) Observe(target string, port int, open bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.obs[target]; !ok {
		m.obs[target] = make(map[int][]bool)
	}
	win := m.obs[target][port]
	win = append(win, open)
	if len(win) > m.cfg.WindowSize {
		win = win[len(win)-m.cfg.WindowSize:]
	}
	m.obs[target][port] = win
}

// Trend returns the current trend for a specific port on a target.
func (m *ScanTrendManager) Trend(target string, port int) PortTrend {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.computeTrend(target, port)
}

// Trends returns all port trends for a target, sorted by port number.
func (m *ScanTrendManager) Trends(target string) []PortTrend {
	m.mu.Lock()
	defer m.mu.Unlock()
	ports, ok := m.obs[target]
	if !ok {
		return nil
	}
	result := make([]PortTrend, 0, len(ports))
	for port := range ports {
		result = append(result, m.computeTrend(target, port))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Port < result[j].Port
	})
	return result
}

func (m *ScanTrendManager) computeTrend(target string, port int) PortTrend {
	win := m.obs[target][port]
	n := len(win)
	t := PortTrend{
		Target:    target,
		Port:      port,
		Direction: TrendStable,
		Samples:   n,
		UpdatedAt: time.Now(),
	}
	if n < m.cfg.MinSamples {
		return t
	}
	var openCount int
	for _, v := range win {
		if v {
			openCount++
		}
	}
	t.OpenRate = float64(openCount) / float64(n)
	if t.OpenRate >= m.cfg.RiseThreshold {
		t.Direction = TrendRising
	} else if (1-t.OpenRate) >= m.cfg.FallThreshold {
		t.Direction = TrendFalling
	}
	return t
}
