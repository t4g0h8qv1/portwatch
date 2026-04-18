package portrank_test

import (
	"testing"

	"github.com/user/portwatch/internal/portrank"
)

func TestScore_CriticalPorts(t *testing.T) {
	r := portrank.New()
	for _, port := range []int{22, 23, 3389} {
		if got := r.Score(port); got != portrank.RankCritical {
			t.Errorf("port %d: want critical, got %s", port, got)
		}
	}
}

func TestScore_HighPorts(t *testing.T) {
	r := portrank.New()
	for _, port := range []int{21, 25, 445, 1433, 3306} {
		if got := r.Score(port); got != portrank.RankHigh {
			t.Errorf("port %d: want high, got %s", port, got)
		}
	}
}

func TestScore_MediumPort(t *testing.T) {
	r := portrank.New()
	if got := r.Score(80); got != portrank.RankMedium {
		t.Errorf("want medium, got %s", got)
	}
}

func TestScore_LowPort(t *testing.T) {
	r := portrank.New()
	if got := r.Score(8080); got != portrank.RankLow {
		t.Errorf("want low, got %s", got)
	}
}

func TestSetOverride(t *testing.T) {
	r := portrank.New()
	r.SetOverride(8080, portrank.RankCritical)
	if got := r.Score(8080); got != portrank.RankCritical {
		t.Errorf("want critical after override, got %s", got)
	}
}

func TestRankAll_SortedDescending(t *testing.T) {
	r := portrank.New()
	ports := []int{8080, 22, 80, 3306}
	results := r.RankAll(ports)
	if len(results) != len(ports) {
		t.Fatalf("want %d results, got %d", len(ports), len(results))
	}
	for i := 1; i < len(results); i++ {
		if results[i].Rank > results[i-1].Rank {
			t.Errorf("results not sorted at index %d", i)
		}
	}
}

func TestRankAll_Empty(t *testing.T) {
	r := portrank.New()
	results := r.RankAll(nil)
	if len(results) != 0 {
		t.Errorf("want empty results, got %d", len(results))
	}
}

func TestRank_String(t *testing.T) {
	cases := map[portrank.Rank]string{
		portrank.RankLow:      "low",
		portrank.RankMedium:   "medium",
		portrank.RankHigh:     "high",
		portrank.RankCritical: "critical",
	}
	for rank, want := range cases {
		if got := rank.String(); got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	}
}
