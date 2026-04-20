package portwatch

import "errors"

var (
	errShadowInvalidObs = errors.New("shadow: minObservations must be >= 1")
	errShadowInvalidAge = errors.New("shadow: maxAge must be positive")
)
