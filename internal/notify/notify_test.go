package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notify"
)

// errNotifier always returns an error.
type errNotifier struct{ msg string }

func (e *errNotifier) Notify(_, _ string) error { return errors.New(e.msg) }

// captureNotifier records calls.
type captureNotifier struct {
	subject string
	body    string
}

func (c *captureNotifier) Notify(subject, body string) error {
	c.subject = subject
	c.body = body
	return nil
}

func TestLogNotifier_WritesOutput(t *testing.T) {
	var buf bytes.Buffer
	n := &notify.LogNotifier{Out: &buf}

	if err := n.Notify("port change", "port 8080 opened"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALERT: port change") {
		t.Errorf("expected subject in output, got: %s", out)
	}
	if !strings.Contains(out, "port 8080 opened") {
		t.Errorf("expected body in output, got: %s", out)
	}
}

func TestMulti_AllCalled(t *testing.T) {
	a := &captureNotifier{}
	b := &captureNotifier{}
	m := notify.NewMulti(a, b)

	if err := m.Notify("subj", "bod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.subject != "subj" || b.subject != "subj" {
		t.Errorf("not all notifiers received subject")
	}
}

func TestMulti_CollectsErrors(t *testing.T) {
	a := &errNotifier{"first error"}
	b := &errNotifier{"second error"}
	m := notify.NewMulti(a, b)

	err := m.Notify("s", "b")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "first error") || !strings.Contains(err.Error(), "second error") {
		t.Errorf("expected both errors in message, got: %v", err)
	}
}

func TestMulti_PartialError(t *testing.T) {
	ok := &captureNotifier{}
	bad := &errNotifier{"boom"}
	m := notify.NewMulti(ok, bad)

	err := m.Notify("s", "b")
	if err == nil {
		t.Fatal("expected error")
	}
	if ok.subject != "s" {
		t.Errorf("successful notifier was not called")
	}
}
