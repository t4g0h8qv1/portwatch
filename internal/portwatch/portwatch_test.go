package portwatch_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/example/portwatch/internal/portwatch"
)

type captureNotifier struct {
	events []portwatch.Event
}

func (c *captureNotifier) Notify(_ context.Context, e portwatch.Event) error {
	c.events = append(c.events, e)
	return nil
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestRun_FirstScanCreatesBaseline(t *testing.T) {
	dir := t.TempDir()
	bl := filepath.Join(dir, "baseline.json")
	cfg := portwatch.Config{
		Host:     "127.0.0.1",
		Ports:    []int{freePort(t)},
		Baseline: bl,
	}
	n := &captureNotifier{}
	if err := portwatch.Run(context.Background(), cfg, n); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(bl); err != nil {
		t.Fatal("baseline file not created")
	}
	if len(n.events) != 0 {
		t.Fatal("notifier should not be called on first run")
	}
}

func TestRun_NoChanges_NotifierNotCalled(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port
	_ = strconv.Itoa(port)

	dir := t.TempDir()
	cfg := portwatch.Config{
		Host:     "127.0.0.1",
		Ports:    []int{port},
		Baseline: filepath.Join(dir, "bl.json"),
	}
	n := &captureNotifier{}
	// First run creates baseline.
	if err := portwatch.Run(context.Background(), cfg, n); err != nil {
		t.Fatal(err)
	}
	// Second run — same ports open.
	if err := portwatch.Run(context.Background(), cfg, n); err != nil {
		t.Fatal(err)
	}
	if len(n.events) != 0 {
		t.Fatalf("expected 0 notifications, got %d", len(n.events))
	}
}

func TestRun_ScanError_ReturnsError(t *testing.T) {
	cfg := portwatch.Config{
		Host:     "256.256.256.256", // invalid
		Ports:    []int{80},
		Baseline: filepath.Join(t.TempDir(), "bl.json"),
	}
	n := &captureNotifier{}
	if err := portwatch.Run(context.Background(), cfg, n); err == nil {
		t.Fatal("expected error for invalid host")
	}
}
