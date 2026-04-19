package portwatch

import (
	"sync"
	"time"
)

// Metrics tracks runtime counters for a portwatch runner.
type Metrics struct {
	mu           sync.Mutex
	Scans        int
	Alerts       int
	Errors       int
	ConsecErrors int
	LastScanAt   time.Time
	LastError    error
}

// RecordScan increments the scan counter and resets consecutive errors.
func (m *Metrics) RecordScan() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Scans++
	m.ConsecErrors = 0
	m.LastScanAt = time.Now()
}

// RecordAlert increments the alert counter.
func (m *Metrics) RecordAlert() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Alerts++
}

// RecordError increments error counters and stores the last error.
func (m *Metrics) RecordError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Errors++
	m.ConsecErrors++
	m.LastError = err
}

// Snapshot returns a copy of the current metrics (safe for reading).
func (m *Metrics) Snapshot() Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	return Metrics{
		Scans:        m.Scans,
		Alerts:       m.Alerts,
		Errors:       m.Errors,
		ConsecErrors: m.ConsecErrors,
		LastScanAt:   m.LastScanAt,
		LastError:    m.LastError,
	}
}
