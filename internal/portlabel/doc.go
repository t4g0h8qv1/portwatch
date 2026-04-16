// Package portlabel resolves port numbers to human-readable service names.
//
// Built-in mappings cover common IANA-registered services. Callers may
// supply a custom map to override or extend the defaults.
//
// Usage:
//
//	l := portlabel.New(nil)
//	fmt.Println(l.Resolve(443))  // "https"
//
//	custom := map[int]string{9000: "myapp"}
//	l2 := portlabel.New(custom)
//	fmt.Println(l2.Resolve(9000)) // "myapp"
package portlabel
