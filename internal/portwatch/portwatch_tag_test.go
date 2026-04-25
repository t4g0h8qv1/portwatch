package portwatch

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewTagManager_Empty(t *testing.T) {
	m := NewTagManager()
	if len(m.Targets()) != 0 {
		t.Fatal("expected no targets")
	}
}

func TestSet_And_Get(t *testing.T) {
	m := NewTagManager()
	_ = m.Set("host1", []string{"prod", "critical"})
	tags := m.Get("host1")
	if len(tags) != 2 || tags[0] != "prod" || tags[1] != "critical" {
		t.Fatalf("unexpected tags: %v", tags)
	}
}

func TestSet_EmptyTarget_ReturnsError(t *testing.T) {
	m := NewTagManager()
	if err := m.Set("", []string{"x"}); err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestSet_DeduplicatesTags(t *testing.T) {
	m := NewTagManager()
	_ = m.Set("host1", []string{"env", "env", "prod"})
	tags := m.Get("host1")
	if len(tags) != 2 {
		t.Fatalf("expected 2 unique tags, got %d: %v", len(tags), tags)
	}
}

func TestGet_Missing_ReturnsNil(t *testing.T) {
	m := NewTagManager()
	if tags := m.Get("unknown"); tags != nil {
		t.Fatalf("expected nil, got %v", tags)
	}
}

func TestRemove_ClearsTags(t *testing.T) {
	m := NewTagManager()
	_ = m.Set("host1", []string{"prod"})
	m.Remove("host1")
	if tags := m.Get("host1"); tags != nil {
		t.Fatalf("expected nil after remove, got %v", tags)
	}
}

func TestTargets_Sorted(t *testing.T) {
	m := NewTagManager()
	_ = m.Set("zzz", []string{"a"})
	_ = m.Set("aaa", []string{"b"})
	_ = m.Set("mmm", []string{"c"})
	targets := m.Targets()
	for i := 1; i < len(targets); i++ {
		if targets[i-1] > targets[i] {
			t.Fatalf("targets not sorted: %v", targets)
		}
	}
}

func TestSet_Overwrites(t *testing.T) {
	m := NewTagManager()
	_ = m.Set("host1", []string{"old"})
	_ = m.Set("host1", []string{"new", "extra"})
	tags := m.Get("host1")
	if len(tags) != 2 || tags[0] != "new" {
		t.Fatalf("unexpected tags after overwrite: %v", tags)
	}
}

func TestWriteTagTable_ContainsHeaders(t *testing.T) {
	m := NewTagManager()
	var buf bytes.Buffer
	WriteTagTable(&buf, m)
	if !strings.Contains(buf.String(), "TARGET") || !strings.Contains(buf.String(), "TAGS") {
		t.Fatalf("missing headers in output: %s", buf.String())
	}
}

func TestWriteTagTable_ShowsTarget(t *testing.T) {
	m := NewTagManager()
	_ = m.Set("myhost", []string{"staging"})
	var buf bytes.Buffer
	WriteTagTable(&buf, m)
	if !strings.Contains(buf.String(), "myhost") {
		t.Fatalf("expected target in output: %s", buf.String())
	}
}
