// Package portexpiry tracks how long individual ports have remained
// continuously open on a scanned host.
//
// A Registry is loaded from (or created at) a JSON file on disk. After each
// scan the caller should call Track with the current set of open ports so that
// newly-opened ports are recorded and closed ports are pruned.
//
// Expired returns all ports whose first-seen age exceeds a caller-supplied
// threshold, allowing portwatch to alert when a port has been open longer than
// policy permits.
//
// Example:
//
//	r, err := portexpiry.Load("/var/lib/portwatch/expiry.json")
//	if err != nil { ... }
//	r.Track(openPorts, time.Now())
//	for _, e := range r.Expired(24*time.Hour, time.Now()) {
//		fmt.Printf("port %d open since %s\n", e.Port, e.FirstSeen)
//	}
//	_ = r.Save()
package portexpiry
