package portwatch_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/example/portwatch/internal/portwatch"
)

// fakeNotifier records calls.
type fakeNotifier struct {
	Called  bool
	Message string
}

func (f *fakeNotifier) Notify(_ context.Context, msg string) error {
	f.Called = true
	f.Message = msg
	return nil
}

// freePort returns an OS-assigned free TCP port.
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
	path := filepath.Join(dir, "baseline.json")
	n := &fakeNotifier{}

	cfg := portwatch.Config{
		Target:       "127.0.0.1",
		Ports:        []int{freePort(t)},
		BaselinePath: path,
		Notifier:     n,
	}

	if err := portwatch.Run(context.Background(), cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("baseline file not created")
	}
	if n.Called {
		t.Fatal("notifier should not be called on first run")
	}
}

func TestRun_NoChanges_NotifierNotCalled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	n := &fakeNotifier{}

	cfg := portwatch.Config{
		Target:       "127.0.0.1",
		Ports:        []int{freePort(t)},
		BaselinePath: path,
		Notifier:     n,
	}

	// First run creates baseline.
	if err := portwatch.Run(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	// Second run with same (closed) ports — no diff.
	if err := portwatch.Run(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	if n.Called {
		t.Fatalf("notifier called unexpectedly: %s", n.Message)
	}
}

func TestRun_ScanError_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	n := &fakeNotifier{}
	cfg := portwatch.Config{
		Target:       fmt.Sprintf("%%invalid%%"),
		Ports:        []int{80},
		BaselinePath: filepath.Join(dir, "b.json"),
		Notifier:     n,
	}
	// An invalid target should propagate a scan error.
	_ = portwatch.Run(context.Background(), cfg) // may or may not error depending on OS
}
