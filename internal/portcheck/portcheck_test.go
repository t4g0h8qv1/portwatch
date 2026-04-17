package portcheck_test

import (
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/portwatch/internal/portcheck"
)

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

func startListener(t *testing.T, port int) net.Listener {
	t.Helper()
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { l.Close() })
	return l
}

func TestRun_CreatesBaselineOnFirstScan(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	port := freePort(t)
	l, _ := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	defer l.Close()

	checker := portcheck.New(path, 200*time.Millisecond)
	res, err := checker.Run("127.0.0.1", []int{port})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Changed() {
		t.Error("first scan should not report changes")
	}
	if _, statErr := os.Stat(path); statErr != nil {
		t.Error("baseline file should exist after first scan")
	}
}

func TestRun_DetectsNewPort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	// First scan — no open ports.
	checker := portcheck.New(path, 200*time.Millisecond)
	_, err := checker.Run("127.0.0.1", []int{})
	if err != nil {
		t.Fatalf("first run: %v", err)
	}

	// Open a port and scan again.
	port := freePort(t)
	l, _ := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	defer l.Close()

	res, err := checker.Run("127.0.0.1", []int{port})
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if !res.Changed() {
		t.Error("expected changes to be detected")
	}
	if len(res.Diff.Opened) != 1 || res.Diff.Opened[0] != port {
		t.Errorf("expected opened port %d, got %v", port, res.Diff.Opened)
	}
}

func TestRun_NoChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	checker := portcheck.New(path, 200*time.Millisecond)
	// Two scans with no open ports.
	for i := 0; i < 2; i++ {
		res, err := checker.Run("127.0.0.1", []int{})
		if err != nil {
			t.Fatalf("run %d: %v", i, err)
		}
		if i > 0 && res.Changed() {
			t.Error("expected no changes on second scan")
		}
	}
}
