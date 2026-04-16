// Package portping measures round-trip latency to open TCP ports.
package portping

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Result holds the outcome of a single ping attempt.
type Result struct {
	Host    string
	Port    int
	Latency time.Duration
	Err     error
}

// Pinger probes TCP ports and measures latency.
type Pinger struct {
	timeout time.Duration
}

// New returns a Pinger with the given per-probe timeout.
func New(timeout time.Duration) (*Pinger, error) {
	if timeout <= 0 {
		return nil, fmt.Errorf("portping: timeout must be positive, got %v", timeout)
	}
	return &Pinger{timeout: timeout}, nil
}

// Ping attempts a TCP connection to host:port and returns latency.
func (p *Pinger) Ping(ctx context.Context, host string, port int) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()

	dialer := &net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	latency := time.Since(start)

	if err != nil {
		return Result{Host: host, Port: port, Err: err}
	}
	conn.Close()
	return Result{Host: host, Port: port, Latency: latency}
}

// PingAll probes each port in ports and returns all results.
func (p *Pinger) PingAll(ctx context.Context, host string, ports []int) []Result {
	results := make([]Result, 0, len(ports))
	for _, port := range ports {
		select {
		case <-ctx.Done():
			results = append(results, Result{Host: host, Port: port, Err: ctx.Err()})
		default:
			results = append(results, p.Ping(ctx, host, port))
		}
	}
	return results
}
