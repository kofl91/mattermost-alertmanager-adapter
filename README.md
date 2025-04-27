
# Mattermost alerts

## E2E Testing

Spin up a preview mattermost instance.
```
docker run --name mattermost-preview -d --publish 8065:8065 mattermost/mattermost-preview
```

Go to http://localhost:8065 and create a new organization by inputing your information.
Once you are in the mattermost organization, create a new channel and then go to integrations -> incoming webhooks -> Add Incoming Webhook and choose your channel. Copy the resulting webhook url.

http://localhost:8065/hooks/mt6f97ska3ggby3tngfkuxeaty

Now you can execute the e2e tests by using:

```
MATTERMOST_WEBHOOK_URL=http://127.0.0.1:8065/hooks/mt6f97ska3ggby3tngfkuxeaty go test ./internal/e2e -v
```

View the alerts channel to see how your alert looks now.