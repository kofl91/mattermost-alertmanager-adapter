package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"mattermost-alertmanager-adapter/internal/adapter"
	"mattermost-alertmanager-adapter/internal/alertmanager"
	"mattermost-alertmanager-adapter/internal/mattermost"
)

func main() {
	mattermostWebhookURL := os.Getenv("MATTERMOST_WEBHOOK_URL")
	if mattermostWebhookURL == "" {
		log.Fatal("MATTERMOST_WEBHOOK_URL environment variable is not set")
		return
	}
	http.HandleFunc("/", createAlertmanagerWebhookHandler(mattermostWebhookURL))
	log.Fatal(http.ListenAndServe(":9997", nil))
}

func createAlertmanagerWebhookHandler(mattermostWebhookURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var alertManagerPayload alertmanager.AlertManagerPayload
		err = json.Unmarshal(body, &alertManagerPayload)
		if err != nil || len(alertManagerPayload.Alerts) == 0 {
			http.Error(w, "Invalid Alertmanager webhook format", http.StatusBadRequest)
			return
		}

		mattermostMessage, err := adapter.MapAlertToMattermostMessage(alertManagerPayload)
		if err != nil {
			http.Error(w, "Failed to convert to Mattermost message", http.StatusInternalServerError)
			return
		}

		err = mattermost.SendToMattermost(mattermostWebhookURL, mattermostMessage)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to send to Mattermost: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Alert received and processed"))
	}
}
