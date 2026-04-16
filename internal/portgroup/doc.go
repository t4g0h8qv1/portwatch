// Package portgroup provides a registry of named port groups for use in
// portwatch configuration and alerting.
//
// Groups allow users to refer to logical service categories (e.g. "web",
// "database") instead of enumerating individual port numbers in every config
// file. Groups are resolved to a deduplicated, ordered list of port numbers
// before scanning.
//
// A set of well-known groups is available via WellKnown():
//
//	reg := portgroup.WellKnown()
//	ports, err := reg.Resolve([]string{"web", "database"})
//
// Custom groups can be registered on top of or instead of the defaults:
//
//	reg := portgroup.New()
//	reg.Register("internal", []int{8080, 9090, 9200})
//	ports, _ := reg.Lookup("internal")
package portgroup
