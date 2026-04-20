// Package portwatch — shadow tracker
//
// ShadowTracker implements a confirmation buffer for newly-discovered open
// ports. Rather than alerting on the very first observation of an unexpected
// port (which can produce noise from transient listeners), the tracker
// requires a port to be seen at least MinObservations times within MaxAge
// before it is considered confirmed and returned to the caller as a real
// finding.
//
// Typical usage:
//
//	st, _ := portwatch.NewShadowTracker(3, 5*time.Minute)
//	if st.Observe("192.168.1.1", 8080) {
//		// port confirmed — raise alert
//	}
//
// Entries that are not confirmed within MaxAge are pruned automatically on
// the next call to Observe or Pending for that target.
package portwatch
