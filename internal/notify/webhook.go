package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// webhookPayload is the JSON body sent to a webhook endpoint.
type webhookPayload struct {
	Timestamp string `json:"timestamp"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

// WebhookNotifier sends alert notifications via HTTP webhook.
type WebhookNotifier struct {
	cfg    WebhookConfig
	client *http.Client
}

// NewWebhookNotifier creates a WebhookNotifier from the given config.
// A zero-value WebhookConfig.Method defaults to POST.
func NewWebhookNotifier(cfg WebhookConfig) *WebhookNotifier {
	if cfg.Method == "" {
		cfg.Method = http.MethodPost
	}
	return &WebhookNotifier{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// Notify sends a JSON payload to the configured webhook URL.
func (w *WebhookNotifier) Notify(subject, body string) error {
	payload := webhookPayload{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Subject:   subject,
		Body:      body,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook marshal: %w", err)
	}

	req, err := http.NewRequest(w.cfg.Method, w.cfg.URL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range w.cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
