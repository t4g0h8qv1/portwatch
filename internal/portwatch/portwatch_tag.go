package portwatch

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"text/tabwriter"
)

// TagManager assigns and retrieves free-form string tags for scan targets.
// Tags are useful for grouping targets by environment, team, or criticality.
type TagManager struct {
	mu   sync.RWMutex
	tags map[string][]string
}

// NewTagManager returns an empty TagManager.
func NewTagManager() *TagManager {
	return &TagManager{
		tags: make(map[string][]string),
	}
}

// Set replaces all tags for the given target. Returns an error if target is
// empty.
func (m *TagManager) Set(target string, tags []string) error {
	if target == "" {
		return fmt.Errorf("portwatch: tag target must not be empty")
	}
	deduped := deduplicateTags(tags)
	m.mu.Lock()
	m.tags[target] = deduped
	m.mu.Unlock()
	return nil
}

// Get returns the tags for the given target. Returns nil if the target has no
// tags.
func (m *TagManager) Get(target string) []string {
	m.mu.RLock()
	v := m.tags[target]
	m.mu.RUnlock()
	out := make([]string, len(v))
	copy(out, v)
	return out
}

// Remove deletes all tags for the given target.
func (m *TagManager) Remove(target string) {
	m.mu.Lock()
	delete(m.tags, target)
	m.mu.Unlock()
}

// Targets returns all targets that have at least one tag, sorted.
func (m *TagManager) Targets() []string {
	m.mu.RLock()
	out := make([]string, 0, len(m.tags))
	for t := range m.tags {
		out = append(out, t)
	}
	m.mu.RUnlock()
	sort.Strings(out)
	return out
}

// WriteTagTable writes a human-readable table of target→tags to w.
func WriteTagTable(w io.Writer, m *TagManager) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tTAGS")
	for _, t := range m.Targets() {
		fmt.Fprintf(tw, "%s\t%v\n", t, m.Get(t))
	}
	tw.Flush()
}

func deduplicateTags(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
