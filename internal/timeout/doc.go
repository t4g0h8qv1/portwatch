// Package timeout provides per-host scan timeout management for portwatch.
//
// A Manager holds a global default timeout and optional per-host overrides.
// Callers obtain the effective timeout via Get, which falls back to the
// default when no override is registered.
//
// Example:
//
//	m, _ := timeout.New(2 * time.Second)
//	m.Set("slow-host.internal", 5 * time.Second)
//	d := m.Get("slow-host.internal") // 5s
//	d  = m.Get("other-host")         // 2s
package timeout
