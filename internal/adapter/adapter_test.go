package adapter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"mattermost-alerts/internal/alertmanager"
	"mattermost-alerts/internal/mattermost"
)

func TestMapAlertToMattermostMessage(t *testing.T) {
	startTime := time.Now().UTC() // Use UTC for predictable time formatting

	// Sample input: one alert
	input := alertmanager.AlertManagerPayload{
		Alerts: []alertmanager.Alert{
			{
				Status: "firing",
				Labels: map[string]string{
					"alertname": "HighCPUUsage",
					"instance":  "server1",
				},
				Annotations: map[string]string{
					"summary": "CPU usage above 90%",
				},
				GeneratorURL: "http://prometheus/graph?g0.expr=cpu",
				StartsAt:     startTime,
			},
		},
	}

	// Expected output
	expected := &mattermost.MattermostMessage{
		Username: "Alertmanager",
		Text:     "Received 1 alert(s)",
		Attachments: []mattermost.Attachment{
			{
				Title: "[firing] HighCPUUsage on server1",
				Text:  "CPU usage above 90%",
				Fields: []mattermost.Field{
					{
						Title: "alertname",
						Value: "HighCPUUsage",
						Short: true,
					},
					{
						Title: "instance",
						Value: "server1",
						Short: true,
					},
					{
						Title: "Links to Reference",
						Value: "[:prometheus:](http://prometheus/graph?g0.expr=cpu) | [:grafana:](http://grafana/dashboard/xyz) | [:alert:](http://alertmanager/alerts) ",
						Short: false,
					},
				},
				Footer: "Started at: " + startTime.Format(time.RFC3339),
			},
		},
	}

	// Call the adapter
	actual, err := MapAlertToMattermostMessage(input)
	assert.NoError(t, err, "expected no error from MapAlertToMattermostMessage")
	assert.Equal(t, expected, actual, "expected and actual Mattermost messages should match")
}
