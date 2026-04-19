package portwatch

import (
	"fmt"
	"sync"
	"time"
)

// PauseManager allows scans to be temporarily paused and resumed.
type PauseManager struct {
	mu      sync.RWMutex
	paused  bool
	until   time.Time
	reason  string
	clock   func() time.Time
}

// NewPauseManager returns a new PauseManager.
func NewPauseManager() *PauseManager {
	return &PauseManager{clock: time.Now}
}

// Pause suspends scanning for the given duration with an optional reason.
func (p *PauseManager) Pause(d time.Duration, reason string) error {
	if d <= 0 {
		return fmt.Errorf("portwatch: pause duration must be positive")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = true
	p.until = p.clock().Add(d)
	p.reason = reason
	return nil
}

// Resume cancels an active pause immediately.
func (p *PauseManager) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = false
	p.until = time.Time{}
	p.reason = ""
}

// IsPaused reports whether scanning is currently paused.
func (p *PauseManager) IsPaused() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.paused {
		return false
	}
	if p.clock().After(p.until) {
		p.paused = false
		p.until = time.Time{}
		p.reason = ""
		return false
	}
	return true
}

// State returns a snapshot of the current pause state.
type PauseState struct {
	Paused bool
	Until  time.Time
	Reason string
}

// State returns the current pause state.
func (p *PauseManager) State() PauseState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return PauseState{
		Paused: p.paused,
		Until:  p.until,
		Reason: p.reason,
	}
}
