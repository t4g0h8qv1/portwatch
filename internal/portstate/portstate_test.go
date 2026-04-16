package portstate_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/portstate"
)

func TestCompare_NoChanges(t *testing.T) {
	prev := portstate.State{Host: "localhost", Ports: []int{22, 80, 443}}
	curr := portstate.State{Host: "localhost", Ports: []int{22, 80, 443}}

	d := portstate.Compare(prev, curr)
	if d.HasChanges() {
		t.Errorf("expected no changes, got opened=%v closed=%v", d.Opened, d.Closed)
	}
}

func TestCompare_NewPorts(t *testing.T) {
	prev := portstate.State{Ports: []int{22}}
	curr := portstate.State{Ports: []int{22, 8080, 9090}}

	d := portstate.Compare(prev, curr)
	if len(d.Opened) != 2 || d.Opened[0] != 8080 || d.Opened[1] != 9090 {
		t.Errorf("unexpected opened ports: %v", d.Opened)
	}
	if len(d.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", d.Closed)
	}
}

func TestCompare_ClosedPorts(t *testing.T) {
	prev := portstate.State{Ports: []int{22, 80, 443}}
	curr := portstate.State{Ports: []int{22}}

	d := portstate.Compare(prev, curr)
	if len(d.Closed) != 2 || d.Closed[0] != 80 || d.Closed[1] != 443 {
		t.Errorf("unexpected closed ports: %v", d.Closed)
	}
	if len(d.Opened) != 0 {
		t.Errorf("expected no opened ports, got %v", d.Opened)
	}
}

func TestCompare_BothChanges(t *testing.T) {
	prev := portstate.State{Ports: []int{22, 80}}
	curr := portstate.State{Ports: []int{22, 443}}

	d := portstate.Compare(prev, curr)
	if !d.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(d.Opened) != 1 || d.Opened[0] != 443 {
		t.Errorf("unexpected opened: %v", d.Opened)
	}
	if len(d.Closed) != 1 || d.Closed[0] != 80 {
		t.Errorf("unexpected closed: %v", d.Closed)
	}
}

func TestCompare_EmptyStates(t *testing.T) {
	d := portstate.Compare(portstate.State{}, portstate.State{})
	if d.HasChanges() {
		t.Error("expected no changes for two empty states")
	}
}
