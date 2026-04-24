package portwatch

import "errors"

var errInvalidDebounceWindow = errors.New("portwatch: debounce window must be positive")
