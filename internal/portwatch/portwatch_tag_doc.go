// Package portwatch provides the core scanning loop and auxiliary managers
// used by the portwatch CLI.
//
// # TagManager
//
// TagManager lets operators attach arbitrary string tags to scan targets so
// that results can be grouped, filtered, or annotated in downstream tooling.
//
// Basic usage:
//
//	m := portwatch.NewTagManager()
//	m.Set("db.internal", []string{"prod", "database"})
//	m.Set("ci.internal",  []string{"staging"})
//
//	tags := m.Get("db.internal") // ["prod", "database"]
//
// Tags are deduplicated on Set.  Removing a target's tags is done with
// Remove.  WriteTagTable renders a human-readable summary to any io.Writer.
package portwatch
