package portwatch_test

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/example/portwatch/internal/portwatch"
)

// fakeNotifier records calls.
type fakeNotifier struct{ called int }

func (f *fakeNotifier) Notify(_ context.Context, _ interface{}) error {
	f.called++
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
	bl := dir + "/baseline.json"
	h := dir + "/history.json"

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port

	cfg := portwatch.Config{
		Target:       "127.0.0.1",
		Ports:        []int{port},
		BaselinePath: bl,
		HistoryPath:  h,
		Timeout:      time.Second,
	}

	res, err := portwatch.Run(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Open) == 0 {
		t.Error("expected at least one open port")
	}
	if _, err := os.Stat(bl); err != nil {
		t.Error("baseline file not created")
	}
	if _, err := os.Stat(h); err != nil {
		t.Error("history file not created")
	}
}

func TestRun_NoChanges_NotifierNotCalled(t *testing.T) {
	dir := t.TempDir()
	bl := dir + "/baseline.json"
	h := dir + "/history.json"

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port

	n := &fakeNotifier{}
	cfg := portwatch.Config{
		Target:       "127.0.0.1",
		Ports:        []int{port},
		BaselinePath: bl,
		HistoryPath:  h,
		Timeout:      time.Second,
		Notifier:     n,
	}

	// first run builds baseline
	if _, err := portwatch.Run(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	// second run should see no changes
	if _, err := portwatch.Run(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	if n.called > 0 {
		t.Errorf("notifier called %d times, want 0", n.called)
	}
}
