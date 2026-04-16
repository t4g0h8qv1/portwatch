package timeout

import "errors"

// ErrHostNotFound is returned when a lookup is performed on an unregistered host
// and strict mode is enabled.
var ErrHostNotFound = errors.New("timeout: no override registered for host")
