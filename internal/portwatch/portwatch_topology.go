package portwatch

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// TopologyEntry records a target's relationship to a named group or region,
// along with the last time the topology was updated.
type TopologyEntry struct {
	Target    string
	Group     string
	Region    string
	UpdatedAt time.Time
}

// TopologyManager tracks logical groupings of scan targets (e.g. region,
// cluster, environment). It is safe for concurrent use.
type TopologyManager struct {
	mu      sync.RWMutex
	entries map[string]TopologyEntry
	now     func() time.Time
}

// NewTopologyManager returns an initialised TopologyManager.
func NewTopologyManager() *TopologyManager {
	return &TopologyManager{
		entries: make(map[string]TopologyEntry),
		now:     time.Now,
	}
}

// Set assigns a group and region to the given target, replacing any previous
// value. It returns an error if target is empty.
func (m *TopologyManager) Set(target, group, region string) error {
	if target == "" {
		return errors.New("portwatch/topology: target must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[target] = TopologyEntry{
		Target:    target,
		Group:     group,
		Region:    region,
		UpdatedAt: m.now(),
	}
	return nil
}

// Get returns the TopologyEntry for the given target and true if it exists,
// or a zero value and false otherwise.
func (m *TopologyManager) Get(target string) (TopologyEntry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[target]
	return e, ok
}

// Remove deletes the topology entry for target. It is a no-op if the target
// is not registered.
func (m *TopologyManager) Remove(target string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, target)
}

// ByGroup returns all targets that belong to the given group, sorted
// alphabetically by target name.
func (m *TopologyManager) ByGroup(group string) []TopologyEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []TopologyEntry
	for _, e := range m.entries {
		if e.Group == group {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Target < out[j].Target })
	return out
}

// ByRegion returns all targets that belong to the given region, sorted
// alphabetically by target name.
func (m *TopologyManager) ByRegion(region string) []TopologyEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []TopologyEntry
	for _, e := range m.entries {
		if e.Region == region {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Target < out[j].Target })
	return out
}

// All returns every registered TopologyEntry sorted alphabetically by target.
func (m *TopologyManager) All() []TopologyEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]TopologyEntry, 0, len(m.entries))
	for _, e := range m.entries {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Target < out[j].Target })
	return out
}
