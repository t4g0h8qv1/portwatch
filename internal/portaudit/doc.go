// Package portaudit correlates live port scan results against a saved
// baseline to produce structured audit records.
//
// Basic usage:
//
//	classifier, _ := severity.New(nil)
//	auditor := portaudit.New(classifier)
//
//	baseline := []int{80, 443}
//	current  := []int{80, 443, 8080}
//
//	rec := auditor.Run("192.168.1.1", baseline, current)
//	if rec.HasChanges() {
//		fmt.Println(rec.Summary())
//	}
//
// The severity of the record reflects the highest-severity new port
// discovered, using the configured severity.Classifier.
package portaudit
