package scanner_test

import (
	"net"
	"testing"
	"time"

	"github.com/example/portwatch/internal/scanner"
)

// startTestListener opens a TCP listener on a random port and returns the port
// number along with a cleanup function.
func startTestListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestScan_OpenPort(t *testing.T) {
	port, cleanup := startTestListener(t)
	defer cleanup()

	opts := scanner.Options{
		Host:    "127.0.0.1",
		Ports:   []int{port},
		Timeout: time.Second,
	}
	results := scanner.Scan(opts)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	opts := scanner.Options{
		Host:    "127.0.0.1",
		Ports:   []int{1},
		Timeout: 300 * time.Millisecond,
	}
	results := scanner.Scan(opts)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestOpenPorts(t *testing.T) {
	results := []scanner.Result{
		{Port: 80, Open: true},
		{Port: 443, Open: false},
		{Port: 8080, Open: true},
	}
	open := scanner.OpenPorts(results)
	if len(open) != 2 {
		t.Fatalf("expected 2 open ports, got %d", len(open))
	}
	if open[0] != 80 || open[1] != 8080 {
		t.Errorf("unexpected open ports: %v", open)
	}
}
