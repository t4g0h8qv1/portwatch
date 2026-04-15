// export_test.go exposes internal fields for white-box testing.
package resolve

import "context"

// SetLookup replaces the DNS lookup function used by the Resolver.
// This is intended for use in tests only.
func (r *Resolver) SetLookup(fn func(ctx context.Context, host string) ([]string, error)) {
	r.lookup = fn
}
