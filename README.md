
# Mattermost alertmanager adapter

This is an adapter between alertmanager and mattermost with the aim to utilise all templating options that mattermost offers and create prettier alerts. This repository was mostly created with AI code generation. Use at your own risk.

I call it matalad because it sounds funny.

## Build and Run

```
go build  -o matalad ./cmd/server/main.go
```

You have to export the MATTERMOST_WEBHOOK_URL environment variable before starting the webserver.
```
export MATTERMOST_WEBHOOK_URL=<URL>
./matalad
```

Example alert message
```
curl -X POST http://localhost:9997/alert \
    -H "Content-Type: application/json" \
    -d '{
    "alerts": [
        {
        "status": "firing",
        "labels": {
            "alertname": "HighCPUUsage",
            "instance": "server1",
            "severity": "critical"
        },
        "annotations": {
            "summary": "CPU usage on server1 is above 90%",
            "description": "The CPU usage on server1 has exceeded the 90% threshold."
        },
        "startsAt": "2025-04-27T15:00:00Z",
        "endsAt": "2025-04-27T15:05:00Z",
        "generatorURL": "http://prometheus.example.com/graph?g0.expr=cpu"
        }
    ]
    }'
```

## Testing

Spin up a preview mattermost instance.
```
docker run --name mattermost-preview -d --publish 8065:8065 mattermost/mattermost-preview
```

Go to http://localhost:8065 and create a new organization by inputing your information.
Once you are in the mattermost organization, create a new channel and then go to integrations -> incoming webhooks -> Add Incoming Webhook and choose your channel. Copy the resulting webhook url.

```
http://localhost:8065/hooks/mt6f97ska3ggby3tngfkuxeaty
```

Now you can execute the e2e tests by using:

```
MATTERMOST_WEBHOOK_URL=http://127.0.0.1:8065/hooks/mt6f97ska3ggby3tngfkuxeaty go test ./internal/e2e -v
```

View the alerts channel to see how your alert looks now.