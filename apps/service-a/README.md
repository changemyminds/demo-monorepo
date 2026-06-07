# service-a (Go + Gin)

A minimal Go HTTP service. Deployed via **Helm** (`deploy/charts/service-a`).

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Liveness, returns `200 {"status":"ok"}` |
| GET | `/version` | Build version (injected via ldflags) |
| GET | `/?name=Ada` | Greeting via the shared `libs/go/common` lib |

## Run locally

```bash
# from repo root (uses go.work for live-at-head libs/go/common)
task run:go APP=service-a
# or:
cd apps/service-a && go run .
```

`PORT` overrides the listen port (default `8080`).

## Test

```bash
task test:go        # tests the whole Go workspace
go test ./...       # from apps/service-a
```

## Build image

```bash
# build context is the REPO ROOT (go.work + shared lib must be present)
task image APP=service-a VERSION=1.2.3
# or:
docker build -f apps/service-a/Dockerfile --build-arg VERSION=1.2.3 -t service-a:1.2.3 .
```
