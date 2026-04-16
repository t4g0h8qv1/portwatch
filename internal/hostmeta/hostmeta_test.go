package hostmeta

import (
	"errors"
	"testing"
	"time"
)

func fixedNow() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestCollect_IPPassthrough(t *testing.T) {
	c := &Collector{
		lookupAddr: func(string) ([]string, error) { return []string{"localhost."}, nil },
		lookupHost: func(string) ([]string, error) { return nil, errors.New("should not call") },
		now:        fixedNow,
	}

	m, err := c.Collect("127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Resolved != "127.0.0.1" {
		t.Errorf("resolved = %q, want 127.0.0.1", m.Resolved)
	}
	if len(m.Hostnames) != 1 || m.Hostnames[0] != "localhost." {
		t.Errorf("hostnames = %v, want [localhost.]", m.Hostnames)
	}
}

func TestCollect_Hostname(t *testing.T) {
	c := &Collector{
		lookupHost: func(host string) ([]string, error) {
			if host == "example.com" {
				return []string{"93.184.216.34"}, nil
			}
			return nil, errors.New("unknown")
		},
		lookupAddr: func(string) ([]string, error) { return []string{"example.com."}, nil },
		now:        fixedNow,
	}

	m, err := c.Collect("example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Resolved != "93.184.216.34" {
		t.Errorf("resolved = %q", m.Resolved)
	}
	if m.Input != "example.com" {
		t.Errorf("input = %q", m.Input)
	}
}

func TestCollect_LookupError(t *testing.T) {
	c := &Collector{
		lookupHost: func(string) ([]string, error) { return nil, errors.New("dns failure") },
		lookupAddr: func(string) ([]string, error) { return nil, nil },
		now:        fixedNow,
	}

	_, err := c.Collect("bad.host")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCollect_ScannedAt(t *testing.T) {
	c := &Collector{
		lookupHost: func(string) ([]string, error) { return []string{"1.2.3.4"}, nil },
		lookupAddr: func(string) ([]string, error) { return nil, nil },
		now:        fixedNow,
	}

	m, err := c.Collect("some.host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.ScannedAt.Equal(fixedNow()) {
		t.Errorf("scanned_at = %v, want %v", m.ScannedAt, fixedNow())
	}
}
