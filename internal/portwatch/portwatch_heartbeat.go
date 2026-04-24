package portwatch

import (
	"io"
	"sync"
	"time"
)

// HeartbeatConfig holds configuration for the heartbeat manager.
type HeartbeatConfig struct {
	// Interval is the expected maximum duration between heartbeats.
	Interval time.Duration
}

// DefaultHeartbeatConfig returns a HeartbeatConfig with sensible defaults.
func DefaultHeartbeatConfig() HeartbeatConfig {
	return HeartbeatConfig{
		Interval: 5 * time.Minute,
	}
}

// heartbeatEntry records the last heartbeat time for a target.
type heartbeatEntry struct {
	last time.Time
}

// HeartbeatManager tracks the last seen heartbeat for each scan target
// and reports whether a target has gone silent.
type HeartbeatManager struct {
	mu      sync.Mutex
	cfg     HeartbeatConfig
	entries map[string]heartbeatEntry
	now     func() time.Time
}

// NewHeartbeatManager creates a HeartbeatManager with the given config.
// Returns an error if the interval is non-positive.
func NewHeartbeatManager(cfg HeartbeatConfig) (*HeartbeatManager, error) {
	if cfg.Interval <= 0 {
		return nil, errInvalidInterval
	}
	return &HeartbeatManager{
		cfg:     cfg,
		entries: make(map[string]heartbeatEntry),
		now:     time.Now,
	}, nil
}

// Beat records a heartbeat for the given target at the current time.
func (h *HeartbeatManager) Beat(target string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries[target] = heartbeatEntry{last: h.now()}
}

// IsSilent reports whether the target has not sent a heartbeat within
// the configured interval. Targets with no recorded heartbeat are silent.
func (h *HeartbeatManager) IsSilent(target string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	e, ok := h.entries[target]
	if !ok {
		return true
	}
	return h.now().Sub(e.last) > h.cfg.Interval
}

// Last returns the time of the last heartbeat for target, and whether
// a heartbeat has been recorded.
func (h *HeartbeatManager) Last(target string) (time.Time, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	e, ok := h.entries[target]
	if !ok {
		return time.Time{}, false
	}
	return e.last, true
}

// WriteHeartbeatTable writes a summary table of heartbeat states to w.
func WriteHeartbeatTable(w io.Writer, h *HeartbeatManager) {
	h.mu.Lock()
	targets := make([]string, 0, len(h.entries))
	for t := range h.entries {
		targets = append(targets, t)
	}
	h.mu.Unlock()

	if len(targets) == 0 {
		io.WriteString(w, "no heartbeat data\n")
		return
	}

	io.WriteString(w, "TARGET\tLAST HEARTBEAT\tSILENT\n")
	for _, t := range targets {
		last, _ := h.Last(t)
		silent := h.IsSilent(t)
		silentStr := "no"
		if silent {
			silentStr = "yes"
		}
		io.WriteString(w, t+"\t"+last.Format(time.RFC3339)+"\t"+silentStr+"\n")
	}
}
