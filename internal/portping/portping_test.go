package portping_test

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portping"
)

func startTCPListener(t *testing.T) (port int, close func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	p, _ := strconv.Atoi(portStr)
	return p, func() { ln.Close() }
}

func TestNew_InvalidTimeout(t *testing.T) {
	_, err := portping.New(0)
	if err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestPing_OpenPort(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	p, _ := portping.New(2 * time.Second)
	res := p.Ping(context.Background(), "127.0.0.1", port)
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Latency <= 0 {
		t.Errorf("expected positive latency, got %v", res.Latency)
	}
}

func TestPing_ClosedPort(t *testing.T) {
	p, _ := portping.New(500 * time.Millisecond)
	res := p.Ping(context.Background(), "127.0.0.1", 1)
	if res.Err == nil {
		t.Fatal("expected error for closed port")
	}
}

func TestPingAll_ReturnsAllResults(t *testing.T) {
	port, stop := startTCPListener(t)
	defer stop()

	p, _ := portping.New(2 * time.Second)
	results := p.PingAll(context.Background(), "127.0.0.1", []int{port, 1})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Err != nil {
		t.Errorf("port %d should be open: %v", port, results[0].Err)
	}
	if results[1].Err == nil {
		t.Errorf("port 1 should be closed")
	}
}

func TestPingAll_ContextCancelled(t *testing.T) {
	p, _ := portping.New(2 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	results := p.PingAll(ctx, "127.0.0.1", []int{80, 443})
	for _, r := range results {
		if r.Err == nil {
			t.Errorf("expected error for cancelled context on port %d", r.Port)
		}
	}
}
