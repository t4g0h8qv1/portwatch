package portmatch_test

import (
	"testing"

	"github.com/example/portwatch/internal/portmatch"
)

func TestAdd_And_Match(t *testing.T) {
	m := portmatch.New()
	if err := m.Add("web", "Web ports", []int{80, 443, 8080}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := m.Match(80)
	if len(got) != 1 || got[0] != "web" {
		t.Errorf("expected [web], got %v", got)
	}
}

func TestMatch_NoRule(t *testing.T) {
	m := portmatch.New()
	_ = m.Add("web", "", []int{80})
	got := m.Match(9999)
	if len(got) != 0 {
		t.Errorf("expected no match, got %v", got)
	}
}

func TestMatch_MultipleRules(t *testing.T) {
	m := portmatch.New()
	_ = m.Add("web", "", []int{80, 443})
	_ = m.Add("plain", "", []int{80})
	got := m.Match(80)
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %v", got)
	}
	if got[0] != "plain" || got[1] != "web" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestMatchAll(t *testing.T) {
	m := portmatch.New()
	_ = m.Add("db", "", []int{5432, 3306})
	result := m.MatchAll([]int{5432, 22, 3306})
	if len(result) != 2 {
		t.Fatalf("expected 2 matched ports, got %d", len(result))
	}
	if _, ok := result[22]; ok {
		t.Error("port 22 should not match")
	}
}

func TestAdd_EmptyName(t *testing.T) {
	m := portmatch.New()
	if err := m.Add("", "", []int{80}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestAdd_NoPorts(t *testing.T) {
	m := portmatch.New()
	if err := m.Add("empty", "", []int{}); err == nil {
		t.Error("expected error for empty port list")
	}
}

func TestAdd_InvalidPort(t *testing.T) {
	m := portmatch.New()
	if err := m.Add("bad", "", []int{0}); err == nil {
		t.Error("expected error for port 0")
	}
	if err := m.Add("bad", "", []int{65536}); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestRules_Sorted(t *testing.T) {
	m := portmatch.New()
	_ = m.Add("zebra", "", []int{1})
	_ = m.Add("alpha", "", []int{2})
	rules := m.Rules()
	if rules[0] != "alpha" || rules[1] != "zebra" {
		t.Errorf("expected sorted rules, got %v", rules)
	}
}
