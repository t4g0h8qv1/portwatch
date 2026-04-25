package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteTrendTable writes a formatted table of port trends to w.
func WriteTrendTable(w io.Writer, trends []PortTrend) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TARGET\tPORT\tDIRECTION\tOPEN RATE\tSAMPLES")
	for _, t := range trends {
		fmt.Fprintf(tw, "%s\t%d\t%s\t%.0f%%\t%d\n",
			t.Target,
			t.Port,
			t.Direction.String(),
			t.OpenRate*100,
			t.Samples,
		)
	}
	tw.Flush()
}

// TrendSummary returns a human-readable summary of trend results.
func TrendSummary(trends []PortTrend) string {
	if len(trends) == 0 {
		return "no trend data"
	}
	var rising, falling, stable int
	for _, t := range trends {
		switch t.Direction {
		case TrendRising:
			rising++
		case TrendFalling:
			falling++
		default:
			stable++
		}
	}
	return fmt.Sprintf("%d ports tracked: %d rising, %d falling, %d stable",
		len(trends), rising, falling, stable)
}
