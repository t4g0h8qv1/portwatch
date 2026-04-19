package portwatch

import (
	"sync"
	"time"
)

// Metrics holds runtime counters for a portwatch run.
type Metrics struct {
	mu sync.Mutex

	ScansTotal    int
	AlertsTotal   int
	ErrorsTotal   int
	LastScanAt    time.Time
	LastAlertAt   time.Time
	OpenPortCount int
}

// RecordScan updates scan-related counters.
func (m *Metrics) RecordScan(openPorts int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ScansTotal++
	m.LastScanAt = time.Now()
	m.OpenPortCount = openPorts
}

// RecordAlert increments the alert counter.
func (m *Metrics) RecordAlert() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AlertsTotal++
	m.LastAlertAt = time.Now()
}

// RecordError increments the error counter.
func (m *Metrics) RecordError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorsTotal++
}

// Snapshot returns a copy of the current metrics.
func (m *Metrics) Snapshot() Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	return Metrics{
		ScansTotal:    m.ScansTotal,
		AlertsTotal:   m.AlertsTotal,
		ErrorsTotal:   m.ErrorsTotal,
		LastScanAt:    m.LastScanAt,
		LastAlertAt:   m.LastAlertAt,
		OpenPortCount: m.OpenPortCount,
	}
}
