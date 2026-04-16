// Package severity provides port change event classification.
//
// A Classifier is configured with sets of critical and warning ports.
// When new or closed ports are detected during a scan, each port is
// assigned a Level (Info, Warning, or Critical) so that downstream
// notifiers and reports can prioritise alerts appropriately.
//
// Example:
//
//	c := severity.New(
//		[]int{22, 3306, 5432},  // critical
//		[]int{80, 8080, 8443},  // warning
//	)
//	level := c.Classify(22) // severity.Critical
package severity
