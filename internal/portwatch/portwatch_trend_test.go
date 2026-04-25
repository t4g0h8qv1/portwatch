package portwatch

import (
	"bytes"
	"strings"
	"testing"
)

func TestDefaultTrendConfig_Defaults(t *testing.T) {
	cfg := DefaultTrendConfig()
	if cfg.WindowSize < 1 {
		t.Fatalf("expected WindowSize >= 1, got %d", cfg.WindowSize)
	}
	if cfg.MinSamples < 1 {
		t.Fatalf("expected MinSamples >= 1, got %d", cfg.MinSamples)
	}
}

func TestNewScanTrendManager_InvalidWindowSize(t *testing.T) {
	cfg := DefaultTrendConfig()
	cfg.WindowSize = 0
	_, err := NewScanTrendManager(cfg)
	if err == nil {
		t.Fatal("expected error for WindowSize=0")
	}
}

func TestNewScanTrendManager_InvalidMinSamples(t *testing.T) {
	cfg := DefaultTrendConfig()
	cfg.MinSamples = 0
	_, err := NewScanTrendManager(cfg)
	if err == nil {
		t.Fatal("expected error for MinSamples=0")
	}
}

func TestNewScanTrendManager_InvalidRiseThreshold(t *testing.T) {
	cfg := DefaultTrendConfig()
	cfg.RiseThreshold = 0
	_, err := NewScanTrendManager(cfg)
	if err == nil {
		t.Fatal("expected error for RiseThreshold=0")
	}
}

func TestNewScanTrendManager_Valid(t *testing.T) {
	cfg := DefaultTrendConfig()
	m, err := NewScanTrendManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
}

func TestObserve_RisingTrend(t *testing.T) {
	cfg := TrendConfig{WindowSize: 4, MinSamples: 3, RiseThreshold: 0.75, FallThreshold: 0.75}
	m, _ := NewScanTrendManager(cfg)
	for i := 0; i < 4; i++ {
		m.Observe("host", 80, true)
	}
	tr := m.Trend("host", 80)
	if tr.Direction != TrendRising {
		t.Errorf("expected rising, got %s", tr.Direction)
	}
}

func TestObserve_FallingTrend(t *testing.T) {
	cfg := TrendConfig{WindowSize: 4, MinSamples: 3, RiseThreshold: 0.75, FallThreshold: 0.75}
	m, _ := NewScanTrendManager(cfg)
	for i := 0; i < 4; i++ {
		m.Observe("host", 443, false)
	}
	tr := m.Trend("host", 443)
	if tr.Direction != TrendFalling {
		t.Errorf("expected falling, got %s", tr.Direction)
	}
}

func TestObserve_InsufficientSamples(t *testing.T) {
	cfg := TrendConfig{WindowSize: 10, MinSamples: 5, RiseThreshold: 0.5, FallThreshold: 0.5}
	m, _ := NewScanTrendManager(cfg)
	m.Observe("host", 22, true)
	m.Observe("host", 22, true)
	tr := m.Trend("host", 22)
	if tr.Direction != TrendStable {
		t.Errorf("expected stable (insufficient samples), got %s", tr.Direction)
	}
}

func TestTrends_SortedByPort(t *testing.T) {
	cfg := DefaultTrendConfig()
	m, _ := NewScanTrendManager(cfg)
	for i := 0; i < cfg.MinSamples; i++ {
		m.Observe("host", 8080, true)
		m.Observe("host", 22, false)
		m.Observe("host", 443, true)
	}
	trends := m.Trends("host")
	for i := 1; i < len(trends); i++ {
		if trends[i].Port < trends[i-1].Port {
			t.Errorf("trends not sorted at index %d", i)
		}
	}
}

func TestTrends_MissingTarget(t *testing.T) {
	cfg := DefaultTrendConfig()
	m, _ := NewScanTrendManager(cfg)
	if got := m.Trends("unknown"); got != nil {
		t.Errorf("expected nil for unknown target, got %v", got)
	}
}

func TestWriteTrendTable_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	WriteTrendTable(&buf, []PortTrend{
		{Target: "h", Port: 80, Direction: TrendRising, OpenRate: 1.0, Samples: 5},
	})
	out := buf.String()
	for _, hdr := range []string{"TARGET", "PORT", "DIRECTION", "OPEN RATE", "SAMPLES"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("missing header %q in output", hdr)
		}
	}
}

func TestTrendSummary_NoTrends(t *testing.T) {
	got := TrendSummary(nil)
	if !strings.Contains(got, "no trend data") {
		t.Errorf("unexpected summary: %q", got)
	}
}

func TestTrendSummary_Counts(t *testing.T) {
	trends := []PortTrend{
		{Direction: TrendRising},
		{Direction: TrendFalling},
		{Direction: TrendStable},
	}
	got := TrendSummary(trends)
	if !strings.Contains(got, "1 rising") || !strings.Contains(got, "1 falling") || !strings.Contains(got, "1 stable") {
		t.Errorf("unexpected summary: %q", got)
	}
}
