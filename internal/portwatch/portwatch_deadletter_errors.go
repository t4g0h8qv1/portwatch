package portwatch

import "errors"

// ErrInvalidDeadLetterSize is returned when maxSize is less than 1.
var ErrInvalidDeadLetterSize = errors.New("portwatch: dead letter queue maxSize must be >= 1")
