package throttle

import "errors"

// ErrInvalidCooldown is returned when a non-positive cooldown duration is provided.
var ErrInvalidCooldown = errors.New("throttle: cooldown must be greater than zero")
