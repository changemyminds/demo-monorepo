# service-b (Go + Gin)

A minimal Go HTTP service. Deployed via **Helm** (`deploy/charts/service-b`).

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Liveness, returns `200 {"status":"ok"}` |
| GET | `/version` | Build version (injected via ldflags) |
| GET | `/?name=Grace` | Greeting via the shared `libs/go/common` lib |

## Run locally

```bash
task run:go APP=service-b
# or:
cd apps/service-b && go run .
```

`PORT` overrides the listen port (default `8080`).

## Test

```bash
task test:go
```

## Build image

```bash
task image APP=service-b VERSION=1.2.3
# or:
docker build -f apps/service-b/Dockerfile --build-arg VERSION=1.2.3 -t service-b:1.2.3 .
```
