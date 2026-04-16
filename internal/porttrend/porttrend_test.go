package porttrend_test

import (
	"testing"

	"github.com/example/portwatch/internal/history"
	"github.com/example/portwatch/internal/porttrend"
)

func makeRecords(portSets [][]int) []history.Record {
	records := make([]history.Record, len(portSets))
	for i, ports := range portSets {
		records[i] = history.Record{OpenPorts: ports}
	}
	return records
}

func TestAnalyze_Empty(t *testing.T) {
	a := porttrend.New(nil)
	if got := a.Analyze(); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestAnalyze_SingleRecord(t *testing.T) {
	records := makeRecords([][]int{{80, 443}})
	a := porttrend.New(records)
	entries := a.Analyze()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Frequency != 1.0 {
			t.Errorf("port %d: expected frequency 1.0, got %f", e.Port, e.Frequency)
		}
	}
}

func TestAnalyze_FrequencySorted(t *testing.T) {
	records := makeRecords([][]int{
		{80, 443, 22},
		{80, 443},
		{80},
	})
	a := porttrend.New(records)
	entries := a.Analyze()

	if entries[0].Port != 80 {
		t.Errorf("expected port 80 first, got %d", entries[0].Port)
	}
	if entries[0].SeenCount != 3 {
		t.Errorf("expected seen 3, got %d", entries[0].SeenCount)
	}
	if entries[0].Frequency != 1.0 {
		t.Errorf("expected freq 1.0, got %f", entries[0].Frequency)
	}
}

func TestUnstable_FiltersLowFrequency(t *testing.T) {
	records := makeRecords([][]int{
		{80, 9999},
		{80},
		{80},
		{80},
	})
	a := porttrend.New(records)
	unstable := a.Unstable(0.5)

	if len(unstable) != 1 {
		t.Fatalf("expected 1 unstable port, got %d", len(unstable))
	}
	if unstable[0].Port != 9999 {
		t.Errorf("expected port 9999, got %d", unstable[0].Port)
	}
}

func TestUnstable_NoneBelow(t *testing.T) {
	records := makeRecords([][]int{{80}, {80}, {80}})
	a := porttrend.New(records)
	if got := a.Unstable(0.5); len(got) != 0 {
		t.Errorf("expected no unstable ports, got %v", got)
	}
}
