package portwatch

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// WriteShardTable writes a summary table of shard assignments to w.
func WriteShardTable(w io.Writer, m *ShardManager) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "SHARD\tTARGETS")
	for i := 0; i < m.ShardCount(); i++ {
		targets, _ := m.Targets(i)
		fmt.Fprintf(tw, "%d\t%d\n", i, len(targets))
	}
	tw.Flush()
}

// ShardSummary returns a one-line summary of shard distribution.
func ShardSummary(m *ShardManager) string {
	total := 0
	max := 0
	for i := 0; i < m.ShardCount(); i++ {
		targets, _ := m.Targets(i)
		n := len(targets)
		total += n
		if n > max {
			max = n
		}
	}
	if total == 0 {
		return fmt.Sprintf("%d shards, no targets assigned", m.ShardCount())
	}
	return fmt.Sprintf("%d shards, %d targets total, max %d per shard", m.ShardCount(), total, max)
}
