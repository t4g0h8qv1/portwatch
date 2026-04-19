package portwatch

import (
	"fmt"
	"io"
	"time"
)

// HealthStatus represents the current health of the portwatch runner.
type HealthStatus int

const (
	HealthOK HealthStatus = iota
	HealthDegraded
	HealthUnknown
)

func (h HealthStatus) String() string {
	switch h {
	case HealthOK:
		return "ok"
	case HealthDegraded:
		return "degraded"
	default:
		return "unknown"
	}
}

// HealthReport holds a point-in-time health snapshot.
type HealthReport struct {
	Status      HealthStatus
	Target      string
	LastScanAt  time.Time
	ConsecErrors int
	LastError   error
}

// Health derives a HealthReport from the current Metrics and config.
func Health(cfg Config, m Metrics) HealthReport {
	r := HealthReport{
		Target:      cfg.Target,
		LastScanAt:  m.LastScanAt,
		ConsecErrors: m.ConsecErrors,
		LastError:   m.LastError,
	}
	switch {
	case m.ConsecErrors == 0:
		r.Status = HealthOK
	case m.ConsecErrors > 0:
		r.Status = HealthDegraded
	default:
		r.Status = HealthUnknown
	}
	return r
}

// WriteHealth writes a human-readable health report to w.
func WriteHealth(w io.Writer, r HealthReport) {
	fmt.Fprintf(w, "target : %s\n", r.Target)
	fmt.Fprintf(w, "status : %s\n", r.Status)
	if r.LastScanAt.IsZero() {
		fmt.Fprintf(w, "last scan : never\n")
	} else {
		fmt.Fprintf(w, "last scan : %s\n", r.LastScanAt.Format(time.RFC3339))
	}
	fmt.Fprintf(w, "consec errors: %d\n", r.ConsecErrors)
	if r.LastError != nil {
		fmt.Fprintf(w, "last error : %s\n", r.LastError)
	}
}
