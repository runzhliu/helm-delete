# helm cm-delete

A Helm plugin to delete chart versions from [ChartMuseum](https://github.com/helm/chartmuseum).

This plugin is the `DELETE` counterpart of [helm cm-push](https://github.com/chartmuseum/helm-push).

## Install

```bash
helm plugin install https://github.com/chartmuseum/helm-cm-delete
```

## Usage

```
helm cm-delete [NAME] [VERSION] [REPO]
```

`REPO` can be a configured repository name (added via `helm repo add`) or a direct URL.

### Examples

```bash
# Delete using a configured repo name
helm cm-delete mychart 1.2.3 myrepo

# Delete using a direct URL
helm cm-delete mychart 1.2.3 https://chartmuseum.example.com

# Delete with basic auth
helm cm-delete mychart 1.2.3 myrepo --username admin --password secret

# Delete using a bearer token
helm cm-delete mychart 1.2.3 myrepo --access-token mytoken

# Delete via a reverse proxy with a context path
helm cm-delete mychart 1.2.3 myrepo --context-path /charts
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--username` | `-u` | Chart repository username |
| `--password` | `-p` | Chart repository password |
| `--access-token` | | Send token in `Authorization: Bearer` header |
| `--auth-header` | | Custom header name for the access token |
| `--context-path` | | URL prefix when ChartMuseum is behind a reverse proxy |
| `--ca-file` | | CA certificate bundle for server verification |
| `--cert-file` | | TLS client certificate |
| `--key-file` | | TLS client private key |
| `--insecure` | `-i` | Skip TLS certificate verification |
| `--timeout` | | Request timeout in seconds (default: 30) |

## Environment variables

The plugin reads the same environment variables as `helm cm-push` for consistency:

| Variable | Description |
|----------|-------------|
| `HELM_REPO_USERNAME` | Repository username |
| `HELM_REPO_PASSWORD` | Repository password |
| `HELM_REPO_ACCESS_TOKEN` | Bearer token |
| `HELM_REPO_AUTH_HEADER` | Custom auth header name |
| `HELM_REPO_CONTEXT_PATH` | Context path prefix |

## ChartMuseum API

The plugin calls `DELETE /api/charts/{name}/{version}` on the target ChartMuseum instance. Make sure `ALLOW_OVERWRITE` or `DISABLE_DELETE` is not set on the server side if you encounter 403 errors.

## Development

```bash
# Build locally
make build

# Run tests
make test

# Install from local source
HELM_CM_DELETE_PLUGIN_NO_INSTALL_HOOK=1 make install
```

## License

Apache 2.0
