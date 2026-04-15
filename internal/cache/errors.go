package cache

import "errors"

// ErrInvalidTTL is returned when a non-positive TTL is provided to New.
var ErrInvalidTTL = errors.New("cache: TTL must be greater than zero")
