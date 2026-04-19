package portwatch

import "errors"

var (
	errDrainTimeout       = errors.New("portwatch: drain timed out waiting for in-flight scans")
	errInvalidDrainTimeout = errors.New("portwatch: drain timeout must be positive")
)
