package filter

import (
	"reflect"
	"sort"
	"testing"
)

func sorted(s []int) []int {
	out := make([]int, len(s))
	copy(out, s)
	sort.Ints(out)
	return out
}

func TestFilter_NoOptions(t *testing.T) {
	ports := []int{22, 80, 443, 8080}
	got := Filter(ports, Options{})
	if !reflect.DeepEqual(sorted(got), sorted(ports)) {
		t.Errorf("expected %v, got %v", ports, got)
	}
}

func TestFilter_AllowList(t *testing.T) {
	ports := []int{22, 80, 443, 8080}
	opts := Options{Allow: []int{80, 443}}
	got := Filter(ports, opts)
	want := []int{80, 443}
	if !reflect.DeepEqual(sorted(got), want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestFilter_DenyList(t *testing.T) {
	ports := []int{22, 80, 443, 8080}
	opts := Options{Deny: []int{22, 8080}}
	got := Filter(ports, opts)
	want := []int{80, 443}
	if !reflect.DeepEqual(sorted(got), want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestFilter_DenyTakesPrecedenceOverAllow(t *testing.T) {
	ports := []int{22, 80, 443}
	opts := Options{
		Allow: []int{22, 80, 443},
		Deny:  []int{22},
	}
	got := Filter(ports, opts)
	want := []int{80, 443}
	if !reflect.DeepEqual(sorted(got), want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	got := Filter([]int{}, Options{Allow: []int{80}, Deny: []int{22}})
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestFilter_AllDenied(t *testing.T) {
	ports := []int{22, 80}
	opts := Options{Deny: []int{22, 80}}
	got := Filter(ports, opts)
	if len(got) != 0 {
		t.Errorf("expected empty slice, got %v", got)
	}
}

func TestFilter_AllowListPortNotInInput(t *testing.T) {
	// Ports in the allow list that are not present in the input should not
	// appear in the output.
	ports := []int{22, 80}
	opts := Options{Allow: []int{80, 443}}
	got := Filter(ports, opts)
	want := []int{80}
	if !reflect.DeepEqual(sorted(got), want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
