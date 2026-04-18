// Package portprofile provides a registry of named port profiles.
//
// A profile is a named collection of ports that describes the expected
// open ports for a particular host role (e.g. "web", "database", "ssh").
//
// Profiles can be registered at runtime or loaded from the built-in
// defaults via Default(). They are used alongside portpolicy and
// portcheck to validate that a scanned host conforms to its declared role.
//
// Example:
//
//	reg := portprofile.Default()
//	profile, err := reg.Get("web")
//	if err != nil { ... }
//	fmt.Println(profile.Ports) // [80 443 8080 8443]
package portprofile
