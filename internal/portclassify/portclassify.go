// Package portclassify categorises ports into broad traffic classes
// (e.g. web, database, remote-access) based on well-known assignments and
// user-supplied overrides.
package portclassify

import "strings"

// Class represents a broad traffic category.
type Class string

const (
	ClassWeb          Class = "web"
	ClassDatabase     Class = "database"
	ClassRemoteAccess Class = "remote-access"
	ClassMail         Class = "mail"
	ClassFile         Class = "file"
	ClassUnknown      Class = "unknown"
)

func (c Class) String() string { return string(c) }

var builtins = map[int]Class{
	21: ClassFile, 22: ClassRemoteAccess, 23: ClassRemoteAccess,
	25: ClassMail, 80: ClassWeb, 110: ClassMail, 143: ClassMail,
	443: ClassWeb, 465: ClassMail, 587: ClassMail, 993: ClassMail,
	995: ClassMail, 3306: ClassDatabase, 5432: ClassDatabase,
	6379: ClassDatabase, 27017: ClassDatabase, 3389: ClassRemoteAccess,
	8080: ClassWeb, 8443: ClassWeb,
}

// Classifier maps ports to traffic classes.
type Classifier struct {
	overrides map[int]Class
}

// New returns a Classifier backed by built-in mappings.
func New() *Classifier {
	return &Classifier{overrides: make(map[int]Class)}
}

// Override registers a custom class for port, replacing any built-in value.
func (c *Classifier) Override(port int, class Class) {
	c.overrides[port] = class
}

// Classify returns the Class for the given port.
func (c *Classifier) Classify(port int) Class {
	if cl, ok := c.overrides[port]; ok {
		return cl
	}
	if cl, ok := builtins[port]; ok {
		return cl
	}
	return ClassUnknown
}

// ClassifyAll returns a map of port → Class for every port in the slice.
func (c *Classifier) ClassifyAll(ports []int) map[int]Class {
	out := make(map[int]Class, len(ports))
	for _, p := range ports {
		out[p] = c.Classify(p)
	}
	return out
}

// ByClass groups the provided ports by their Class.
func (c *Classifier) ByClass(ports []int) map[Class][]int {
	out := make(map[Class][]int)
	for _, p := range ports {
		cl := c.Classify(p)
		out[cl] = append(out[cl], p)
	}
	return out
}

// ParseClass converts a string to a Class, returning ClassUnknown if
// unrecognised.
func ParseClass(s string) Class {
	switch Class(strings.ToLower(s)) {
	case ClassWeb, ClassDatabase, ClassRemoteAccess, ClassMail, ClassFile:
		return Class(strings.ToLower(s))
	}
	return ClassUnknown
}
