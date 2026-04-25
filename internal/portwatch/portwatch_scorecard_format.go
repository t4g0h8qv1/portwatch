package portwatch

import "fmt"

// ScorecardSummary returns a one-line human-readable summary of scorecard
// entries suitable for log output or status lines.
func ScorecardSummary(entries []ScorecardEntry) string {
	if len(entries) == 0 {
		return "scorecard: no targets tracked"
	}
	var totalScans, totalAlerts, totalErrors int
	for _, e := range entries {
		totalScans += e.TotalScans
		totalAlerts += e.TotalAlerts
		totalErrors += e.TotalErrors
	}
	return fmt.Sprintf(
		"scorecard: %d target(s), %d scan(s), %d alert(s), %d error(s)",
		len(entries), totalScans, totalAlerts, totalErrors,
	)
}

// AlertRate returns the fraction of scans that produced an alert for the given
// entry. It returns 0 when no scans have been recorded.
func AlertRate(e ScorecardEntry) float64 {
	if e.TotalScans == 0 {
		return 0
	}
	return float64(e.TotalAlerts) / float64(e.TotalScans)
}

// ErrorRate returns the fraction of scans that resulted in an error for the
// given entry. It returns 0 when no scans have been recorded.
func ErrorRate(e ScorecardEntry) float64 {
	if e.TotalScans == 0 {
		return 0
	}
	return float64(e.TotalErrors) / float64(e.TotalScans)
}
