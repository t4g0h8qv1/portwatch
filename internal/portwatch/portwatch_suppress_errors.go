package portwatch

import "errors"

var errInvalidSuppressTTL = errors.New("portwatch: suppress TTL must be greater than zero")
