package portwatch

import (
	"fmt"
	"sort"
	"sync"
)

// LabelManager assigns and retrieves human-readable labels for scan targets.
// Labels are used in reports and notifications to provide context beyond raw
// hostnames or IP addresses.
type LabelManager struct {
	mu     sync.RWMutex
	labels map[string]string
}

// NewLabelManager returns an initialised LabelManager.
func NewLabelManager() *LabelManager {
	return &LabelManager{
		labels: make(map[string]string),
	}
}

// Set assigns a label to a target. Both target and label must be non-empty.
func (m *LabelManager) Set(target, label string) error {
	if target == "" {
		return fmt.Errorf("portwatch: label target must not be empty")
	}
	if label == "" {
		return fmt.Errorf("portwatch: label must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.labels[target] = label
	return nil
}

// Get returns the label for target, or an empty string if none is set.
func (m *LabelManager) Get(target string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.labels[target]
}

// Label returns the label for target if set, otherwise the target itself.
func (m *LabelManager) Label(target string) string {
	if l := m.Get(target); l != "" {
		return l
	}
	return target
}

// Remove deletes the label for target. It is a no-op if target has no label.
func (m *LabelManager) Remove(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.labels, target)
}

// All returns a sorted slice of all (target, label) pairs.
func (m *LabelManager) All() [][2]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([][2]string, 0, len(m.labels))
	for t, l := range m.labels {
		out = append(out, [2]string{t, l})
	}
	sort.Slice(out, func(i, j int) bool { return out[i][0] < out[j][0] })
	return out
}
