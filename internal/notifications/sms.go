package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SMSSender struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

func NewSMSSender(sid, token, from string) *SMSSender {
	return &SMSSender{
		AccountSID: sid,
		AuthToken:  token,
		FromNumber: from,
	}
}

func (s *SMSSender) Send(to, message string) error {
	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.AccountSID)

	payload := map[string]string{
		"To":   to,
		"From": s.FromNumber,
		"Body": message,
	}

	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.SetBasicAuth(s.AccountSID, s.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("Twilio error: %s", resp.Status)
	}

	fmt.Println("📱 SMS sent to:", to)
	return nil
}
