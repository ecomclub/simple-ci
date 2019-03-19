# simple-ci
Go HTTP service to deploy Node.js apps (Systemd services) from GitHub

## Environment variables

Configure the service with following env variables:

Variable              | Sample
---                   | ---
`LOG_FILE`            | `/var/log/app.out`
`APPS_ROOT`           | `/var/apps/`
`SERVICES_PREFIX`     | `node-`
`TCP_PORT`            | `:30000`

## Endpoint

```http
http://localhost:3000/?AppName=myapp&GitBranch=master
```

## GitHub setup

The simplest way to get it running is creating a **webhook** on GitHub
repository after releases (recommended) or pushes.

### Secret

You should use a proxy server (Nginx) and check the `X-Hub-Signature`
header sent by GitHub for a specific secret hash.
