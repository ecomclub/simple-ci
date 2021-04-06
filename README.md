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
http://localhost:3000/?AppName=myapp&GitBranch=master&Secret=123hash
```

### Proxy

You should use a proxy server (_Nginx_) and handle this request
from a simpler URL, eg.:

```http
https://myserver.com/myapp/deploy
```

```conf
# nginx.conf
location ~ ^\/(?<app>[a-z0-9-]+)\/deploy$ {
  proxy_pass http://[::1]:3000/?Secret=SuperSecretPassword&AppName=$app&GitBranch=master;
}
```

## GitHub setup

The simplest way to get it running is creating a **webhook** on GitHub
repository after releases (recommended) or pushes.
