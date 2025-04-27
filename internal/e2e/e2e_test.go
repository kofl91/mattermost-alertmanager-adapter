package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"mattermost-alertmanager-adapter/internal/adapter"
	"mattermost-alertmanager-adapter/internal/alertmanager"
	"mattermost-alertmanager-adapter/internal/mattermost"
)

// loadWebhookURL loads the webhook URL from env variable
func loadWebhookURL(t *testing.T) string {
	url := os.Getenv("MATTERMOST_WEBHOOK_URL")
	require.NotEmpty(t, url, "MATTERMOST_WEBHOOK_URL must be set for E2E tests")
	return url
}

// sendToMattermost posts the MattermostMessage to the webhook
func sendToMattermost(t *testing.T, webhookURL string, msg *mattermost.MattermostMessage) {
	bodyBytes, err := json.Marshal(msg)
	require.NoError(t, err, "failed to marshal MattermostMessage")

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(bodyBytes))
	require.NoError(t, err, "failed to POST to Mattermost webhook")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "expected HTTP 200 OK from Mattermost")
}

func TestEndToEndAlertToMattermost(t *testing.T) {
	webhookURL := loadWebhookURL(t)

	// Build a fake Alertmanager payload
	startTime := time.Now().UTC()

	alertPayload := alertmanager.AlertManagerPayload{
		Alerts: []alertmanager.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "HighMemoryUsage",
					"instance":  "server2",
				},
				Annotations: map[string]string{
					"summary": "Memory usage above 90%",
				},
				GeneratorURL: "http://prometheus/graph?g0.expr=mem",
				StartsAt:     startTime,
			},
		},
	}

	// Convert it to Mattermost message
	mmMessage, err := adapter.MapAlertToMattermostMessage(alertPayload)
	require.NoError(t, err, "failed to map alert to Mattermost message")

	// Send to Mattermost webhook
	sendToMattermost(t, webhookURL, mmMessage)
}
