# Go Prometheus sentry exporter

based on the [prometheus sentry exporter](https://github.com/ujamii/prometheus-sentry-exporter) but written in go.

# Arguments

| Flag              | Default | Description             |
|-------------------| ------- |-------------------------|
| addr              | :8080 | http service address    |
| sentry-url        | https://sentry.io | sentry url              |
| sentry-org        | n/a | Sentry ORG              |
| sentry-auth-token | n/a | Sentry Auth token       |
| sentry-base | api/0 | sentry base |


## Sentry Auth token 

You can generate this from https://sentry.io/settings/{organization_slug}/developer-settings/).

### Full steps

1. Go to https://sentry.io.
2. Navigate to Organization Settings, and then to Developer Settings.
3. Select New Internal Integration.
3. Use a valid name such as prometheus.
4. Go to PERMISSIONS, provide Read permissions to the required resources such as "Project", "Issue and Event", and "Organization".
5. Copy the token, and then use this token when configuring the data source within Grafana.

```Note: In Sentry, The Admin, Manager, or Owner role is required to get an internal integration token```


