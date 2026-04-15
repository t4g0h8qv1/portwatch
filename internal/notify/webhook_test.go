package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestWebhookNotifier_Success(t *testing.T) {
	var received map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(notify.WebhookConfig{URL: ts.URL})
	if err := n.Notify("test subject", "test body"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received["subject"] != "test subject" {
		t.Errorf("expected subject 'test subject', got %q", received["subject"])
	}
	if received["body"] != "test body" {
		t.Errorf("expected body 'test body', got %q", received["body"])
	}
	if received["timestamp"] == "" {
		t.Errorf("expected non-empty timestamp")
	}
}

func TestWebhookNotifier_CustomHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Token") != "secret" {
			t.Errorf("expected X-Token header")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(notify.WebhookConfig{
		URL:     ts.URL,
		Headers: map[string]string{"X-Token": "secret"},
	})
	if err := n.Notify("s", "b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWebhookNotifier_Non2xx(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := notify.NewWebhookNotifier(notify.WebhookConfig{URL: ts.URL})
	err := n.Notify("s", "b")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestWebhookNotifier_BadURL(t *testing.T) {
	n := notify.NewWebhookNotifier(notify.WebhookConfig{URL: "http://127.0.0.1:1"})
	if err := n.Notify("s", "b"); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
