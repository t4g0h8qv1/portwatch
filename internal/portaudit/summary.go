package portaudit

import "fmt"

// Summary returns a human-readable one-line description of the record.
func (r Record) Summary() string {
	if !r.HasChanges() {
		return fmt.Sprintf("%s: no changes (%d stable)", r.Host, len(r.Stable))
	}
	return fmt.Sprintf("%s: +%d new, -%d gone, %d stable [%s]",
		r.Host, len(r.New), len(r.Gone), len(r.Stable), r.Severity)
}
