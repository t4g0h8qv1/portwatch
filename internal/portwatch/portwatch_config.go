package portwatch

import "time"

// Config holds runtime options for a portwatch Run.
type Config struct {
	// Target is the host to scan (hostname or IP).
	Target string

	// Ports is the list of port numbers to probe.
	Ports []int

	// Timeout is the per-port dial timeout.
	Timeout time.Duration

	// BaselinePath is the file path used to persist the port baseline.
	BaselinePath string

	// Interval controls how often the scheduler triggers a scan.
	// A zero value means run once and exit.
	Interval time.Duration

	// AlertOnGone, when true, fires notifications when previously open
	// ports are no longer detected.
	AlertOnGone bool
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Timeout:      500 * time.Millisecond,
		BaselinePath: "portwatch_baseline.json",
		AlertOnGone:  false,
	}
}

// Validate returns an error if the Config is not usable.
func (c Config) Validate() error {
	if c.Target == "" {
		return ErrNoTarget
	}
	if len(c.Ports) == 0 {
		return ErrNoPorts
	}
	if c.Timeout <= 0 {
		return ErrInvalidTimeout
	}
	return nil
}
