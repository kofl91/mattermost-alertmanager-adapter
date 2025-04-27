package adapter

import (
	"fmt"
	"strings"
	"time"

	"mattermost-alertmanager-adapter/internal/alertmanager"
	"mattermost-alertmanager-adapter/internal/mattermost"
)

// MapAlertToMattermostMessage takes an Alertmanager alert and creates a Mattermost message payload.
func MapAlertToMattermostMessage(payload alertmanager.AlertManagerPayload) (*mattermost.MattermostMessage, error) {
	if len(payload.Alerts) == 0 {
		return nil, fmt.Errorf("no alerts found in payload")
	}

	var attachments []mattermost.Attachment

	for _, alert := range payload.Alerts {
		title := createTitleFromAlert(alert)

		text := createTextFromAlert(alert)

		// var actions []mattermost.Action

		// actions = appendPrometheusQueryButton(alert, actions)
		// actions = appendSilenceButton(alert, actions)

		urls := []struct {
			emoji string
			url   string
		}{
			{emoji: ":prometheus:", url: "http://prometheus/graph?g0.expr=cpu"},
			{emoji: ":grafana:", url: "http://grafana/dashboard/xyz"},
			{emoji: ":alert:", url: "http://alertmanager/alerts"},
		}

		attachments = append(attachments, mattermost.Attachment{
			Title:  title,
			Text:   text,
			Fields: append(mapLabelsToFields(alert.Labels), mapLinksToField(urls)),
			Footer: fmt.Sprintf("Started at: %s", alert.StartsAt.Format(time.RFC3339)),
		})
	}

	// Compose the final Mattermost message
	message := &mattermost.MattermostMessage{
		Username:    "Alertmanager",
		Text:        fmt.Sprintf("Received %d alert(s)", len(payload.Alerts)),
		Attachments: attachments,
	}

	return message, nil
}

func appendSilenceButton(alert alertmanager.Alert, actions []mattermost.Action) []mattermost.Action {
	if silenceURL, ok := alert.Annotations["silenceURL"]; ok {
		actions = append(actions, mattermost.Action{
			Id:   "silence_alert",
			Name: "Silence Alert",
			Type: "button",
			Integration: &mattermost.Integration{
				URL: silenceURL,
			},
		})
	}
	return actions
}

func appendPrometheusQueryButton(alert alertmanager.Alert, actions []mattermost.Action) []mattermost.Action {
	if alert.GeneratorURL != "" {
		actions = append(actions, mattermost.Action{
			Id:   "view_in_prometheus",
			Name: "View in Prometheus",
			Type: "button",
			Integration: &mattermost.Integration{
				URL:     alert.GeneratorURL,
				Context: map[string]interface{}{},
			},
		})
	}
	return actions
}

func createTextFromAlert(alert alertmanager.Alert) string {
	text := alert.Annotations["summary"]
	if text == "" {
		text = alert.Annotations["description"]
	}
	if text == "" {
		text = "No summary provided."
	}
	return text
}

func createTitleFromAlert(alert alertmanager.Alert) string {
	alertName := alert.Labels["alertname"]
	instance := alert.Labels["instance"]
	status := alert.Status

	title := fmt.Sprintf("[%s] %s on %s", status, alertName, instance)
	return title
}

// Helper function to map labels into Mattermost fields (key-value pairs).
func mapLabelsToFields(labels map[string]string) []mattermost.Field {
	var fields []mattermost.Field
	for k, v := range labels {
		// Hide noisy internal labels if needed
		if strings.HasPrefix(k, "__") {
			continue
		}
		fields = append(fields, mattermost.Field{
			Title: k,
			Value: v,
			Short: true,
		})
	}
	return fields
}

// mapLinksToField will create a field containing clickable emoji links
func mapLinksToField(urls []struct {
	emoji string
	url   string
}) mattermost.Field {
	// Build the field value with clickable emojis
	var emojiLinks string
	for _, link := range urls {
		emojiLinks += fmt.Sprintf("[%s](%s) | ", link.emoji, link.url)
	}

	// Trim the last "| " from the string
	emojiLinks = emojiLinks[:len(emojiLinks)-2]

	// Return the field with the links
	return mattermost.Field{
		Title: "Links to Reference",
		Value: emojiLinks,
		Short: false,
	}
}
