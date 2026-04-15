// Package portrange provides utilities for parsing port range expressions
// such as "80", "1-1024", or "22,80,443,8000-9000".
package portrange

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	MinPort = 1
	MaxPort = 65535
)

// Parse converts a port range expression into a sorted, deduplicated slice of
// port numbers. Supported formats: "80", "1-1024", "22,80,443,8000-9000".
func Parse(expr string) ([]int, error) {
	seen := make(map[int]struct{})
	var ports []int

	for _, part := range strings.Split(expr, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, err := parsePort(bounds[0])
			if err != nil {
				return nil, fmt.Errorf("invalid range start %q: %w", bounds[0], err)
			}
			hi, err := parsePort(bounds[1])
			if err != nil {
				return nil, fmt.Errorf("invalid range end %q: %w", bounds[1], err)
			}
			if lo > hi {
				return nil, fmt.Errorf("range start %d is greater than end %d", lo, hi)
			}
			for p := lo; p <= hi; p++ {
				if _, dup := seen[p]; !dup {
					seen[p] = struct{}{}
					ports = append(ports, p)
				}
			}
		} else {
			p, err := parsePort(part)
			if err != nil {
				return nil, err
			}
			if _, dup := seen[p]; !dup {
				seen[p] = struct{}{}
				ports = append(ports, p)
			}
		}
	}
	return ports, nil
}

func parsePort(s string) (int, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("not a number: %q", s)
	}
	if n < MinPort || n > MaxPort {
		return 0, fmt.Errorf("port %d out of range [%d, %d]", n, MinPort, MaxPort)
	}
	return n, nil
}
