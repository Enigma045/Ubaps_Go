package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	SMSServer = "https://app.sms8.io"
	SMSApiKey = "0f1b46dfdcccfa1d8461de539375a98bb162f909"
)

// SendSMS sends a single SMS message using the sms8.io gateway
func SendSMS(number, message string) error {
	endpoint := fmt.Sprintf("%s/services/send.php", SMSServer)

	data := url.Values{}
	data.Set("number", number)
	data.Set("message", message)
	data.Set("key", SMSApiKey)
	data.Set("devices", "0") // Use default/primary device
	data.Set("type", "sms")

	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read SMS response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMS gateway returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			Messages []interface{} `json:"messages"`
		} `json:"data"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		// If it's not JSON, it might be a raw error message
		return fmt.Errorf("failed to parse SMS response: %s", string(body))
	}

	if !result.Success {
		return fmt.Errorf("SMS gateway error: %s", result.Error.Message)
	}

	return nil
}
