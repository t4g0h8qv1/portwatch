package portwatch_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/example/portwatch/internal/portwatch"
)

// fakeNotifier records calls.
type fakeNotifier struct {
	called  int
	lastMsg string
}

func (f *fakeNotifier) Notify(_ context.Context, msg string) error {
	f.called++
	f.lastMsg = msg
	return nil
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestRun_FirstScanCreatesBaseline(t *testing.T) {
	dir := t.TempDir()
	n := &fakeNotifier{}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port

	cfg := portwatch.Config{
		Target:       "127.0.0.1",
		Ports:        []int{port},
		BaselinePath: dir + "/baseline.json",
		Timeout:      500 * time.Millisecond,
		Notifier:     n,
	}

	if err := portwatch.Run(context.Background(), cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.called != 0 {
		t.Errorf("expected no notification on first scan, got %d", n.called)
	}
}

func TestRun_NoChanges_NotifierNotCalled(t *testing.T) {
	dir := t.TempDir()
	n := &fakeNotifier{}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port

	cfg := portwatch.Config{
		Target:       "127.0.0.1",
		Ports:        []int{port},
		BaselinePath: dir + "/baseline.json",
		Timeout:      500 * time.Millisecond,
		Notifier:     n,
	}

	// First run — creates baseline.
	if err := portwatch.Run(context.Background(), cfg); err != nil {
		t.Fatalf("first run: %v", err)
	}
	// Second run — same ports, no diff.
	if err := portwatch.Run(context.Background(), cfg); err != nil {
		t.Fatalf("second run: %v", err)
	}
	if n.called != 0 {
		t.Errorf("expected no notification, got %d", n.called)
	}
}
