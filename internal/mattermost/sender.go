// sendToMattermost sends the Mattermost message to the specified webhook URL
package mattermost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendToMattermost(webhookURL string, message *MattermostMessage) error {
	messageBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(messageBody))
	if err != nil {
		return fmt.Errorf("failed to send message to Mattermost: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message to Mattermost: received status %s", resp.Status)
	}

	return nil
}
