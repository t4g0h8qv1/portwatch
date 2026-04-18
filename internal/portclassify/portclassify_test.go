package portclassify_test

import (
	"testing"

	"github.com/user/portwatch/internal/portclassify"
)

func TestClassify_BuiltIn(t *testing.T) {
	c := portclassify.New()
	cases := []struct {
		port int
		want portclassify.Class
	}{
		{80, portclassify.ClassWeb},
		{443, portclassify.ClassWeb},
		{22, portclassify.ClassRemoteAccess},
		{3306, portclassify.ClassDatabase},
		{5432, portclassify.ClassDatabase},
		{25, portclassify.ClassMail},
		{21, portclassify.ClassFile},
		{9999, portclassify.ClassUnknown},
	}
	for _, tc := range cases {
		if got := c.Classify(tc.port); got != tc.want {
			t.Errorf("Classify(%d) = %q, want %q", tc.port, got, tc.want)
		}
	}
}

func TestClassify_Override(t *testing.T) {
	c := portclassify.New()
	c.Override(9200, portclassify.ClassDatabase)
	if got := c.Classify(9200); got != portclassify.ClassDatabase {
		t.Fatalf("expected database, got %q", got)
	}
	// override beats built-in
	c.Override(80, portclassify.ClassDatabase)
	if got := c.Classify(80); got != portclassify.ClassDatabase {
		t.Fatalf("override did not take precedence, got %q", got)
	}
}

func TestClassifyAll(t *testing.T) {
	c := portclassify.New()
	result := c.ClassifyAll([]int{80, 443, 3306})
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[3306] != portclassify.ClassDatabase {
		t.Errorf("expected database for 3306, got %q", result[3306])
	}
}

func TestByClass(t *testing.T) {
	c := portclassify.New()
	groups := c.ByClass([]int{80, 443, 22, 3306, 9999})
	if len(groups[portclassify.ClassWeb]) != 2 {
		t.Errorf("expected 2 web ports, got %v", groups[portclassify.ClassWeb])
	}
	if len(groups[portclassify.ClassUnknown]) != 1 {
		t.Errorf("expected 1 unknown port, got %v", groups[portclassify.ClassUnknown])
	}
}

func TestParseClass_Valid(t *testing.T) {
	for _, s := range []string{"web", "database", "remote-access", "mail", "file"} {
		if got := portclassify.ParseClass(s); got.String() != s {
			t.Errorf("ParseClass(%q) = %q", s, got)
		}
	}
}

func TestParseClass_Invalid(t *testing.T) {
	if got := portclassify.ParseClass("nope"); got != portclassify.ClassUnknown {
		t.Errorf("expected unknown, got %q", got)
	}
}
