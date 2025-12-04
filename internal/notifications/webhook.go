package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WebhookPayload struct {
	Event     string      `json:"event"`
	UserID    string      `json:"user_id"`
	EntityID  string      `json:"entity_id"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type WebhookSender struct{}

func NewWebhookSender() *WebhookSender {
	return &WebhookSender{}
}

func (s *WebhookSender) Send(url string, payload WebhookPayload) error {
	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook rejected: %s", resp.Status)
	}

	fmt.Println("🔗 Webhook delivered to:", url)
	return nil
}
