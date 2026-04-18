// Package portmatch provides pattern-based matching of port numbers against
// named rule sets.
//
// A Matcher holds a collection of Rules, each associating a name and
// description with a set of port numbers. Given one or more open ports, the
// matcher returns which rules apply, enabling human-readable categorisation of
// scan results (e.g. "web", "database", "ssh").
//
// Usage:
//
//	m := portmatch.Default()
//	names := m.Match(443)          // ["web"]
//	all  := m.MatchAll(openPorts)  // map[port][]ruleName
package portmatch
