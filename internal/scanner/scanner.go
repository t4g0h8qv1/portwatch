package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Result holds the result of a port scan for a single port.
type Result struct {
	Port  int
	Open  bool
	Error error
}

// Options configures the port scanner.
type Options struct {
	Host        string
	Ports       []int
	Timeout     time.Duration
	Concurrency int
}

// Scan scans the given ports on the host and returns a slice of Results.
func Scan(opts Options) []Result {
	if opts.Timeout == 0 {
		opts.Timeout = 2 * time.Second
	}
	if opts.Concurrency <= 0 {
		opts.Concurrency = 100
	}

	sem := make(chan struct{}, opts.Concurrency)
	results := make([]Result, len(opts.Ports))
	var wg sync.WaitGroup

	for i, port := range opts.Ports {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx, p int) {
			defer wg.Done()
			defer func() { <-sem }()
			results[idx] = probePort(opts.Host, p, opts.Timeout)
		}(i, port)
	}

	wg.Wait()
	return results
}

// OpenPorts filters a Results slice and returns only open ports.
func OpenPorts(results []Result) []int {
	var open []int
	for _, r := range results {
		if r.Open {
			open = append(open, r.Port)
		}
	}
	return open
}

func probePort(host string, port int, timeout time.Duration) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return Result{Port: port, Open: false, Error: err}
	}
	conn.Close()
	return Result{Port: port, Open: true}
}
