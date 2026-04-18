// Package portrank assigns a risk rank to open ports.
//
// Ports are scored using built-in heuristics (well-known dangerous services
// receive higher ranks) and optional per-port overrides set by the caller.
//
// Usage:
//
//	r := portrank.New()
//	r.SetOverride(8443, portrank.RankHigh)
//
//	results := r.RankAll(openPorts)
//	for _, res := range results {
//		fmt.Printf("port %d — %s\n", res.Port, res.Rank)
//	}
package portrank
