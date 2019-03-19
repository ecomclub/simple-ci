# simple-ci
Go HTTP service to deploy Node.js apps (Systemd services) from GitHub

# GitHub setup

The simplest way to get it running is creating a webhook on GitHub
repository after releases or pushes.

## Secret

You should use a proxy server (Nginx) and check the `X-Hub-Signature`
header sent by GitHub for a specific secret hash.
