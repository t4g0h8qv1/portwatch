package portwatch

import (
	"sync"
	"time"
)

// DrainConfig controls how the drain manager waits for in-flight scans.
type DrainConfig struct {
	Timeout time.Duration
	PollInterval time.Duration
}

// DefaultDrainConfig returns sensible defaults.
func DefaultDrainConfig() DrainConfig {
	return DrainConfig{
		Timeout:      30 * time.Second,
		PollInterval: 100 * time.Millisecond,
	}
}

// DrainManager tracks in-flight scans and allows callers to wait until all
// active scans complete before shutdown.
type DrainManager struct {
	mu     sync.Mutex
	active map[string]int
	cfg    DrainConfig
}

// NewDrainManager returns a DrainManager with the given config.
func NewDrainManager(cfg DrainConfig) (*DrainManager, error) {
	if cfg.Timeout <= 0 {
		return nil, errInvalidDrainTimeout
	}
	if cfg.PollInterval <= 0 {
		cfg.PollInterval = DefaultDrainConfig().PollInterval
	}
	return &DrainManager{active: make(map[string]int), cfg: cfg}, nil
}

// Acquire marks a scan as started for the given target.
func (d *DrainManager) Acquire(target string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.active[target]++
}

// Release marks a scan as finished for the given target.
func (d *DrainManager) Release(target string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.active[target] > 0 {
		d.active[target]--
	}
	if d.active[target] == 0 {
		delete(d.active, target)
	}
}

// InFlight returns the total number of active scans across all targets.
func (d *DrainManager) InFlight() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	total := 0
	for _, n := range d.active {
		total += n
	}
	return total
}

// Wait blocks until all in-flight scans complete or the configured timeout
// elapses. Returns errDrainTimeout if the deadline is exceeded.
func (d *DrainManager) Wait() error {
	deadline := time.Now().Add(d.cfg.Timeout)
	for {
		if d.InFlight() == 0 {
			return nil
		}
		if time.Now().After(deadline) {
			return errDrainTimeout
		}
		time.Sleep(d.cfg.PollInterval)
	}
}
